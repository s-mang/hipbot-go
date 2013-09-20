Botling
=====

#### A Hipchat bot witten in Go[lang] 

Botling is a neat little bot with some awesome functionality. He sits comfortably in your Hipchat room and obeys your every request (the ones he's familiar with anyway). Botling knows how to search for nearby restaurants, get an image given a tag, search the New York Times, get a weather forecast, more! (see below for full command list)

#### Go Knowledge Optional!

### [Setup & Configuration Instructions](https://github.com/Sproutling/botling/wiki)

### Botling Command Examples
1. `@botling nearby sushi`
    * Google Places search for "sushi" near your specified lat & lng
2. `@botling nytimes technology`
    * Search recent New York Times articles in the "technology" section
3. `@botling image me sunset`
    * Search Flickr for a photo of a "sunset"
4. `@botling weather me tomorrow`
    * Get tomorrow's weather
5. `@botling trivia me today`
    * Get trivia factoids for today's date (interesting or historical events)
6. `@botling trivia me 123`
    * Get trivia on a number "123" (math facts)
7. `@botling wolfram me periodic table of elements`
    * Get computational information from Wolfram Alpha on the "periodic table of elements"
8. `@botling gopkg math`
    * Get the synopsis and path for the Go[lang] package "math"
9. `@botling company logo`
    * Get the logo of your company (url specified by you - see configuration section of the wiki)
10. `@botling thesaurus me challenge`
    * Get synonyms for the word "challenge"
11. `@botling search me HMAC`
    * Search the web for information on "HMAC" (uses duckduckgo.com for search engine)
12. `@botling goodnight`
    * Tell Botling goodnight! He will respond with a nice farewell
13. `@botling foobar baz`
    * Say anything else, and botling will respond with "Hello, FirstName LastName"
