# GoH
Utility lib for writing extremely simple webhooks in go, among other things.

# Packages

For now, GoH has four packages.

1. `gohclient`: http helper methods to easily communicate with a REST api written in go.

2. `gohcmd`: methods to ease out the proper creation of cmd utilities.

3. `gohserver`: http helper methods to make it easy to create webhooks.

4. `gohtypes`: helper types for handling with webhook constructs

# Examples

You can find good examples that use goh in the following repositories:

[Bindman DNS Webhook](http://github.com/labbsr0x/bindman-dns-webhook/)

[Bindman DNS Swarm Listener](https://github.com/labbsr0x/bindman-dns-swarm-listener)

[Bindman DNS Bind9 Manager](https://github.com/labbsr0x/bindman-dns-bind9)