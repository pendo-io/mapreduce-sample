Mapreduce Sample Application
============================

This is a simple sample application for [the pendo-io/mapreduce go implementation.][1] It works in in both the production and development environments. For either one, you need to place the mapreduce implementation somewhere the sample app can get to it. Here's one way:

    # export GOPATH=$PWD
    # mkdir src
    # goapp get github.com/pendo-io/mapreduce

The application counts the number of times each word occurs in the Gutenburg Project's version of Pride and Prejudice, which has been split into five pieces in the `testdata/` directory ("word" is used loosely here; they're really space separated pieces; punctuation marks are included in words). This is a classic map/reduce sample problem, and it's easy to find an explanation of how it's implemented on the Internet (the documentation for Google's python map/reduce library would be one good place to start). Correct output is in `testdata/pandp-results` for comparison (it was computed using a simple python script to provide a valid comparison).

Running on the Development Server
--------------------------------------

To run the app against the local development server, run:

    # dev_appserver.py .

Then connect to `http://localhost:8080/run`. At the login screen, be sure to click the checkbox to log in as an administrator. This will run the map/reduce process, and provide a link to the status. You can refresh the page until the job is complete, when it will also display links to the output files in the blobstore.

  [1]: www.github.com/pendo-io/mapreduce
  
Running on Appengine
--------------------

Upload the sample application to appengine using `goapp deploy` or `appcfg.py`, perhaps like this:

    # appcfg.py -A mytestapp -V 1 update .
    
Once it's deployed, connect to the `/run` URI for the project and version (and perhaps module) you chose. If you're using the default version of the primary module, something like `https://mytestapp.appspot.com/run` will work. You must use an administrator's account for this url.

The url will tell you the job has been started and provide a link to the job's status. You can manually refresh the link, and when the job is complete it will present links to the files in blobstore so you can download the results.