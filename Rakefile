namespace :backend do
    task :build do
        puts "Building app-backend"
        
        sh "mkdir", "-p", "#{__dir__}/build"
        sh "go", "build", "-o", "#{__dir__}/build/app-backend", "#{__dir__}/cmd/app-backend"
    end

    task :run => :build do
        puts "Running app-backend"

        sh "#{__dir__}/build/app-backend" 
    end
end