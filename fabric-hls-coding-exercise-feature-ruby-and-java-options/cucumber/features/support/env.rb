require 'net/http'
require 'fileutils'
require 'childprocess'
require 'tempfile'
require 'httparty'
require 'json_spec'
require 'erubis'
require 'aws-sdk'
require 'json'
require 'pry-byebug'
require 'colorize'

World(JsonSpec::Helpers, JsonSpec::Matchers)
JsonSpec.configure do
  exclude_keys 'created_at', 'updated_at'
end

LOG_DIR = 'logs'
APPLICATION_LOG_FILE = File.join(LOG_DIR, 'application.log')
HEALTHCHECK_ENDPOINT = 'http://127.0.0.1:1337/healthcheck'
CUCUMBER_BASE = '.'

def wait_till_up_or_timeout
  healthy = false
  i = 0
  puts 'Waiting for application under test to start...'
  while !healthy && i < 30

    unless @app_process.alive?
      shutdown
      puts "The application failed to start up successfully. Check #{APPLICATION_LOG_FILE} for details".red
      exit 1
    end

    begin
      response = Net::HTTP.get_response(URI.parse(HEALTHCHECK_ENDPOINT))
      if response.code == '200'
        healthy = true
      else
        puts 'Health check returned status code: ' + response.code.light_black
      end
    rescue StandardError => e
      puts "GET #{HEALTHCHECK_ENDPOINT} -- #{e}".light_black
    end
    i += 1
    sleep(1) unless healthy
  end

  unless healthy
    shutdown
    puts 'Application failed to pass healthchecks within an acceptable amount of time. Declining to run tests.'.red
    exit 1
  end
end

def startup
  puts 'Starting application...'
  @app_process = ChildProcess.build('../run.sh')
  @app_process.io.stdout = File.new(APPLICATION_LOG_FILE, 'w')
  @app_process.io.stderr = @app_process.io.stdout
  @app_process.leader = true
  @app_process.start
end

def shutdown
  @app_process.stop
end

startup
wait_till_up_or_timeout

at_exit do
  shutdown
end
