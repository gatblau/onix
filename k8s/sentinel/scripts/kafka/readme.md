<img src="../../pics/sentinel_small.png" align="right" height="200" width="200"/>

# Kafka development environment

In order to have Apache Kafka running locally for development purposes.

Ensure [Vagrant](https://www.vagrantup.com/) is installed on the development machine as described [here](https://www.vagrantup.com/docs/installation/).

Then, execute the [go.sh](go.sh) script from this directory.

- The script downloads and launches a [Vagrant box](Vagrantfile).
- Then copies the [setup.sh script](setup.sh) to it and executes it.
- In turn, [setup.sh](setup.sh) downloads and launches Apache Kafka and then creates a topic for Sentinel to send messages to.
- The Virtual Machine exposes zookeeper and kafka ports (2181 and 9092 respectively) on the host running the Virtual Machine.

```bash
# install, configure and run kafka
sh go.sh

# ssh into the kafka box
vagrant ssh

# start the message consumer for the k8s topic
sh kafka-console-consumer.sh --bootstrap-server localhost:9092 --topic k8s --from-beginning
```

[[back to index](../readme.md)]

[*] _The Sentinel icon was made by [Freepik](https://www.freepik.com) from [Flaticon](https://www.flaticon.com) and is licensed by [Creative Commons BY 3.0](http://creativecommons.org/licenses/by/3.0)_