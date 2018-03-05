# Homunculus

How many times have you found yourself telling your Ops staff '_If condition X occurs Y times, you'll need to restart my CF app..._'?  Wouldn't it be awesome if you could tell that to your app instead?  

# This is that awesome
This library is designed to allow you to set arbitrary conditions and thresholds at startup.  When you catch these errors at runtime (e.g. DB connection errors, HTTP timeouts), you simply increment the condition.  Once the condition is breached, Homunculus will restart _the problematic instance_, not the entire microservice deployment.  It will attempt to use a shared Redis cache to allow you a centralized monitoring point; barring that, it will use an in-memory k:v map to track the current state.