# Background and Purpose

This isn't your typical subdomain enumeration tool. This tool has the ability to expand your reach by finding subdomains owned by a common organization, then dive down into the details to find those that are both live and match to your scope.

First, it checks to see if the domain entered as a command line argument is enrolled as a tenant with Microsoft Defender for Identity (MDI). If found in MDI, it returns a list of related domains. Next, it uses the Project Discovery Chaos API to return all subdomains. Finally, it checks to see if the domain resolves to a live IP address. If it resolves, the domain, IP address, and netblock owner are printed to the terminal. This is helpful for extracting domains which are in scope for your pentest.

This is far more effective than simply checking crt.sh, because if the domain is found in MDI, then all domains in that MDI tenant are owned by the same company. Also, I've had much better results querying the Chaos API than when using crt.sh.

# Setup

## Install libxml2-utils:

```
sudo apt install libxml2-utils
```

## Install the Chaos client: (Requires Golang)

```
go install -v github.com/projectdiscovery/chaos-client/cmd/chaos@latest
```

## Set the Chaos client API key:

Signup to get your API key here: https://chaos.projectdiscovery.io/docs/quick-start

Set the key in your .bashrc or .zshrc file:

```
export CHAOS_KEY=<your API key here>
```

# Execution

```
chmod +x dnsrecon.sh
./dnsrecon.sh <domain> 
```

Please be patient. The script doesn't implement parallel execution, so large domains may take a while before you start seeing output.

There are no banners or progress bars. The output is semi-colon separated for easy parsing with awk/grep/cut/Excel.
