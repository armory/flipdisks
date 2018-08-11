# Armory Flipdisks Project
Here at Armory, we have an art project using flipdisk to display messages
and other content. This project is structured as a monorepo, with 3 main projects.
- Web Admin client
- Hardware Controller
- Flipdisk Server


# Equipment List
#### AlphaZeta
Armory has negotiated a 3% discount for you w/ AlfaZeta on flip disc panel orders of 50 pieces or more.
Just mention "Armory Discount" when you order when you email Marcin at info@AlfaZeta.pl, they're also available at +48.42.689.1200.

#### Raspberry Pi
This is our main hardware controller for the boards.  
[Non-affiliate link](https://www.raspberrypi.org/products/raspberry-pi-3-model-b/)


#### USB to RS485 Serial Data Converter
A simple USB data converter. 1 of these can talk to 10 (14x28) AlphaZeta panels.  
[Non-affiliate link](https://www.amazon.com/gp/product/B0721BB8PQ)



# Deploying controller
After working on the code, you can deploy it by:
```bash
ssh-copy-id pi@192.168.86.26  # add your key for easy login

./bin/build-controller.sh
./bin/deploy-controller.sh
```


# Tips and Tricks
### flipdisk-controller deamon
To check on the status on the service, you can do
```bash
journalctl -f -u flipdisk.service
```
