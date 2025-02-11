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

# Example output

```
vhnokia-gw.example.com;12.32.90.147;AT&T Services, Inc.
vhvpn.example.com;12.32.90.179;AT&T Services, Inc.
view.example.com;63.85.196.11;Verizon Business
viewlab.uk.example.com;91.240.17.234;
visionapp-den.services.example.com;64.73.81.63;Example Technologies LLC
visionapp-msp.services.example.com;64.73.65.63;Example Technologies LLC
visionapp.services.example.com;64.73.65.63;Example Technologies LLC
vmrc.uk.example.com;194.105.149.148;
voicetotext.example.com;52.149.215.200;Microsoft Corporation
voip.example.com;12.32.91.66;AT&T Services, Inc.
webapppd.azr.example.com;52.162.107.30;Microsoft Corporation
webapps.example.com;52.149.215.200;Microsoft Corporation
webappsstg.example.com;52.149.215.200;Microsoft Corporation
webcon15.ms.example.com;64.73.23.241;Example Technologies LLC
webcon15.s3.example.com;64.73.23.241;Example Technologies LLC
webmail.example.com;52.149.215.200;Microsoft Corporation
webobjects.example.com;184.24.66.246;Akamai Technologies, Inc.
webobjects2.example.com;184.24.66.246;Akamai Technologies, Inc.
webobjects2beta.example.com;96.7.27.87;Akamai Technologies, Inc.
webportal.example.com;12.32.90.94;AT&T Services, Inc.
webstage.example.com;184.25.36.105;Akamai Technologies, Inc.
webstage.wip.example.com;64.73.16.247;Example Technologies LLC
wifienroll.example.com;52.149.215.200;Microsoft Corporation
www.example.com;184.24.66.246;Akamai Technologies, Inc.
www.mpt.example.com;137.135.86.148;Microsoft Corp
www.uk.example.com;184.24.64.58;Akamai Technologies, Inc.
xenmobilelab.uk.example.com;91.240.17.10;
12-32-90-100.example.com;12.32.90.100;AT&T Services, Inc.
12-32-90-101.example.com;12.32.90.101;AT&T Services, Inc.
```
