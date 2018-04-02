#!/bin/bash
test -f "/bin/bash" && echo "This system has a bash shell"
echo APP_ENV = $APP_ENV
if [ ${APP_ENV} = production ];
then
    app;
else
	go get github.com/pilu/fresh &&
	fresh;
fi
