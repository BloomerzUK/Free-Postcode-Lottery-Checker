# Free-Postcode-Lottery-Checker
This is a little tool i quickly hacked together to send me results of the Free Postcode Lotter (http://freepostcodelottery.com).
Uses Mandrill for email notifications and Rollbar for error handling. 

It's built to run as a service in the background:
free-postcode-lottery-checker.exe --service=install
free-postcode-lottery-checker.exe --service=uninstall
free-postcode-lottery-checker.exe --service=start
free-postcode-lottery-checker.exe --service=stop
