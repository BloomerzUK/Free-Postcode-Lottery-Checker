# Note
This is now super old and doesn't work with current FPL. It's not even worth looking at, as nothing here is relevant anymore. 
Keeping it just for myself for reference, as I might upgrade it one day.




# Free-Postcode-Lottery-Checker
This is a little tool i quickly hacked together to send me results of the Free Postcode Lotter (http://freepostcodelottery.com).
Uses Mandrill for email notifications and Rollbar for error handling. 



It's built to run as a service in the background:

free-postcode-lottery-checker.exe --service=install

free-postcode-lottery-checker.exe --service=uninstall

free-postcode-lottery-checker.exe --service=start

free-postcode-lottery-checker.exe --service=stop
