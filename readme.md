When working with Wordpress, its always best to develop locally and then deploy. But that often means that you want to replicate the content from the live site on your local machine. So you create an export, but there could easily be thousands of attachments that you might need to import to get that imported locally which is going to take a while. Then there may be hundreds of posts that you don't need on the local server. 

This is the beginning of a tool that will cleanup your WXP (Wordpress Export file). It starts by removing all the attachments you aren't using. Then you can delete additional posts by entering the post ids, or a range of post ids. 

Its still a work in progress, but it helped me get the Datadog export from 1000 attachments down to 300 and then further to 100, going from 100 posts to 17. Works pretty well. x
