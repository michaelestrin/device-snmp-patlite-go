Device service for SNMP Patlite (in Go)
=======================================

Requisite
---------
- core-data
- core-metadata
- core-command

GET calls
---------
Get allows you to get the state of the lights and Buzzer.  Here are the allowed GET commands.
- Issuing GET command to: http://device service address:49992/api/v1/device/<device id>/RedLight
- Issuing GET command to: http://device service address:49992/api/v1/device/<device id>/GreenLight
- Issuing GET command to: http://device service address:49992/api/v1/device/<device id>/AmberLight
- Issuing GET command to: http://device service address:49992/api/v1/device/<device id>/Buzzer
- Issuing GET command to: http://device service address:49992/api/v1/device/<device id>/AllLights


PUT calls
---------
Issuing PUT commands to the same URLs above (with PUT method) to set the state of the lights or buzzer.  Put commands require params like following for the Red Light.  
[{"RedLightControlState":"2"}, {"RedLightTimer":"0"} ]

Setting the lights or buzzer requires also providing the timer value.  The timer is in a numer of seconds before executing the state change.

Get allows you to get the state of the lights and Buzzer.  The lights and buzzer state are represented by integers with the following values:

**PATLITE light values**
- LIGHT_OFF   = 1
- LIGHT_ON    = 2
- LIGHT_BLINK = 3
- LIGHT_FLASH = 5

**PATLITE buzzer values**
- BUZZ_ON       = 5
- BUZZ_OFF      = 1
- BUZZ_PATTERN1 = 2
- BUZZ_PATTERN2 = 3
- BUZZ_PATTERN3 = 4

SNMP Tools
==========
For debugging of SNMP commands, install SNMP deamon:  
sudo apt-get install snmpd

sudo apt-get install snmp snmp-mibs-downloader

Some SNMP commands:

snmpset  -v2c -cprivate 192.168.0.20 1.3.6.1.4.1.20440.4.1.5.1.2.1.2.3 i 1 1.3.6.1.4.1.20440.4.1.5.1.2.1.3.3 i 0

snmpget -v2c -c public 192.168.0.20 1.3.6.1.4.1.20440.4.1.5.1.2.1.4.1

snmpwalk -mALL -v1 -cpublic 192.168.0.20 system

snmptest 192.168.0.14

snmptable -v1 localhost


