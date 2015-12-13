require 'json'

class IntegrationTester
  def initialize(verbose=false)
    @alerter = Alerter.new.tap { |a| a.verbose = verbose }
  end

  def execute
    puts "Resolving existing alerts..."
    @alerter.resolve_all_alerts
    alerts = @alerter.get_alerts
    raise "should have no alerts" unless alerts.empty?

    puts "Adding new alert..."
    @alerter.add_alert
    alerts = @alerter.get_alerts
    raise "should have 1 alert" unless alerts.keys.length == 1

    puts "Tests PASSED"
  end
end

class Alerter
  attr_accessor :verbose

  def initialize(deployment_id='987654321')
    @deployment_id = deployment_id
  end

  def get_alerts
    curl_it('GET', "deployments/#{@deployment_id}")
  end

  def add_alert
    post_alert(1, "132435465768798#{rand(5)}")
  end

  def resolve_all_alerts
    alerts = get_alerts
    alerts.each_pair do |capsule_id, alert|
      post_alert(0, capsule_id)
    end
  end

  private

  def post_alert(status, capsule_id)
    body = <<-TEXT
    {
      "client": "localhost",
      "check": {
        "name": "redis0-redis_role",
        "capsule_name": "redis0",
        "output": "no master found",
        "status": #{status},
        "capsule_id": "#{capsule_id}",
        "deployment_id": "#{@deployment_id}",
        "account": "compose"
      }
    }
    TEXT

    curl_it('POST', 'alerts', body)
  end

  def curl_it(*args)
    Curl.new(*args).tap { |c| c.verbose = @verbose }.execute
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
