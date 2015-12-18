require 'json'

class IntegrationTester
  def initialize(verbose=false)
    @verbose = verbose
  end

  def execute
    @deployment = Deployment.new.tap { |d| d.verbose = @verbose }
    @deployment.settings = notification_settings
    puts 'Testing with deployment settings...'
    test_deployment

    @deployment.settings = {}
    @deployment.group_settings = notification_settings
    puts 'Testing with group settings...'
    test_deployment
  end

  def test_deployment
    puts 'Adding deployment...'
    setup_deployment

    puts 'Cleaning up/resolving existing alerts...'
    @alerter = Alerter.new(@deployment.id).tap { |a| a.verbose = @verbose }
    @alerter.resolve_all_alerts
    alerts = @alerter.get_alerts
    raise 'should have no alerts' unless alerts.empty?

    puts 'Adding new alert...'
    @alerter.add_alert_for_check_name(@deployment.check.name)
    alerts = @alerter.get_alerts
    raise 'should have 1 alert' unless alerts.keys.length == 1

    puts 'Resolving new alert...'
    @alerter.resolve_all_alerts
    alerts = @alerter.get_alerts
    raise 'should have no alerts' unless alerts.empty?

    puts 'Tests PASSED'
  end

  def setup_deployment
    @deployment.add_deployment
  end

  def notification_settings
    {'pagerduty_key' => '545bf39a778d45b1b4160a7fd782fae9'} # free, trial account
  end
end

class Curler
  attr_accessor :verbose

  def curl_it(*args)
    Curl.new(*args).tap { |c| c.verbose = @verbose }.execute
  end

  def raise_on_fail
    result = yield
    raise "Failed call: #{result}" unless result.keys.join == 'ok'
  end
end

class Group < Curler
  attr_accessor :group_id, :settings

  def initialize(group_id)
    @group_id = group_id
    @settings = {}
  end

  def add_group
    raise_on_fail { post_group }
  end

  private

  def post_group
    body = {
      'id' => @group_id,
      'settings' => @settings
    }.to_json

    curl_it('POST', 'groups', body)
  end
end

class Deployment < Curler
  attr_accessor :id, :group_id, :settings, :check, :group_settings

  def initialize
    @id = rand(100000).to_s
    @group_id = nil
    @type = 'foobar'
    @settings = {}
    @group_settings = {}
  end

  def add_deployment
    group = Group.new(gid).tap { |g| g.verbose = @verbose; g.settings = @group_settings }
    group.add_group

    raise_on_fail { post_deployment }

    @check = Check.new.tap { |c| c.type = @type; c.verbose = @verbose }
    @check.add_check
  end

  private

  def gid
    @group_id || "g#{@id}"
  end

  def post_deployment
    body = {
      'id' => @id,
      'group_id' => gid,
      'type' => @type,
      'name' => 'foobar-is-awesome',
      'settings' => @settings
    }.to_json

    curl_it('POST', 'deployments', body)
  end
end

class Check < Curler
  attr_accessor :name, :type, :level, :title, :description

  def initialize
    @name = 'foobar_check'
    @type = 'foobar'
    @level = 1
    @title = 'Check the foobar'
    @description = 'Owning a foobar requires checking it.'
  end

  def add_check
    raise_on_fail { post_check }
  end

  private

  def post_check
    body = {
      'name' => @name,
      'type' => @type,
      'level' => @level,
      'title' => @title,
      'description' => @description
    }.to_json

    curl_it('POST', 'checks', body)
  end
end

class Alerter < Curler
  attr_accessor :deployment_id

  def initialize(deployment_id='987654321')
    @deployment_id = deployment_id
    @type = 'redis_role'
    @capsule_name = 'redis0'
  end

  def get_alerts
    curl_it('GET', "deployments/#{@deployment_id}")
  end

  def add_alert_for_check_name(type)
    @type = type
    @status = 1
    @capsule_id = "132435465768798#{rand(5)}"
    post_alert
  end

  def resolve_all_alerts
    alerts = get_alerts
    alerts.each_pair do |capsule_id, alert|
      @status = 0
      @capsule_id = capsule_id
      post_alert
    end
  end

  private

  def post_alert
    body = {
      'client' => 'localhost',
      'check' => {
        'name' => "#{@capsule_name}-#{@type}",
        'capsule_name' => @capsule_name,
        'output' => 'no master found',
        'status' => @status,
        'capsule_id' => @capsule_id,
        'deployment_id' => @deployment_id,
        'account' => 'compose'
      }
    }.to_json

    curl_it('POST', 'alerts', body)
  end
end

# Why shell out to curl instead of net/http or a curl gem or any number of other
# Ruby solutions (httparty, etc)? I like to include an exploratory CLI for
# API apps that clients could pull and would dump actual curl commands to the
# console to more generically demonstrate how to use the API. This isn't fully
# exploratory, but covers some integration basics while doing some refactoring.
class Curl
  attr_accessor :body, :verbose

  def initialize(verb, path, body='')
    @verb = verb
    @path = path
    @body = body
  end

  def execute
    arg_body = @body.empty? ? '' : "-d '#{@body}'"
    cmd = "curl -s -u x:$COMPOSE_SERVICE_PASSWORD -X#{@verb} 'http://localhost:8000/margo/#{@path}' #{arg_body}"
    puts cmd if @verbose
    res = `#{cmd}`
    system "echo '#{res}' | jsonpp" if @verbose # there's a gem for this if you wanna
    JSON.parse(res)
  end
end

IntegrationTester.new(ARGV.include?('-v') || ARGV.include?('--verbose'))
  .execute
