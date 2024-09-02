# PPP

A Healthy Product in a Healthy Business - Prevent, Profile, Perfect

### Logging

- No trails, no tales (No logs, no way to know what happened, or maybe it never happened)
- This could have been a debug (Not everything has to be an info or an error, more noisy logs that are useful can be at debug level)
- When everything is an error, nothing is really an error (Not everything that is logged at error level is really an error, many are just simple warnings, use warn level)
- Too much info, is no info (info logs are useful, but too many of them increases the noise to signal ratio, use them wisely)
- Keep it DRY (Logging the same thing at multiple layers in the code is redundant and noisy)
- No context, no content (Logging a message without the context of what the message is about, is not useful and adds to the noise)
- No traces, no clues (Logging is useful, but in a microservices architecture, it is important to log the trace id to correlate logs across services)
- BE SENSITIVE, NOT EXCLUSIVE
