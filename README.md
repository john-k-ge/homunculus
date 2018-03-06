# Homunculus

How many times have you found yourself telling your Ops staff _"If condition X occurs Y times, you'll need to restart my CF app..."_? 
Instead, don't you wish there was some awesome thing you could tell instead?  Just declare a set conditions and the awesome 
thing would keep score and restart when necessary? 


# This is that awesome thing
This library is designed to allow you to set conditions and thresholds at startup.  When you catch these errors at runtime 
(and we all handle our errors, correct?), you simply increment the condition.  Once the condition is breached, Homunculus 
will restart _the problematic instance_, not the entire microservice deployment.  Ideally, it will attempt to use a shared `predix-cache` 
instance to provide a centralized monitoring point; barring that, it will use an in-memory k:v store to track the current state.

# Getting started
In the `sample` directory, you'll find a simple web app.  You'll need to update the manifest with your Predix org robot 
credentials as well as your current CF API and UAA hostnames (those for US-West are already provided).  
Once you've provided that, you can build and deploy the app with the simple `build_and_deploy.sh`. Once deployed, it will 
register a pair of conditions: `db-error` and `connection-timeout`. You can `curl https://yourapp...predix.io/db` until 
it dies and respawns.  You'll notice that the instances will die off independently, as each instance has a separate tally 
maintained in the cache.  
In a real app, these would represent conditions your app would encounter at runtime: RabbitMQ, Redis, and/or Postgres connection 
failures.  If you trap these errors like you should, you can also increment the associated counter.  Once the app has hit 
the ceiling you've set, Homunculus will throw in the towel and ask Cloud Foundry to restart this app instance. 


