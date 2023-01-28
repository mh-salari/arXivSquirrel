# Welcome to arXivSquirrel <img src="https://arxiv.ai-hue.ir/resources/arXivSquirrel.png" alt="logo" width="100"/>

arXivSquirrel is a Go program that generates personalized RSS feeds from [arXiv.org](https://arxiv.org/) based on user keywords. With arXivSquirrel, you can easily keep track of the latest papers on your areas of interest without having to sift through the overwhelming number of new papers shared on arXiv every day.

One of the unique features of arXivSquirrel is that it includes the image version of the first five pages of each paper in the RSS feed. This allows you to quickly get a sense of the content of the paper without having to read the full title or abstract. This feature can save you a lot of time and help you identify papers that you may want to read in more detail.


## How to use arXivSquirrel

1. Download the arXivSquirrel program from this GitHub repository.

2. Change the `keywords.csv` file based on the keywords you are interested in

3. Open the `main.go` file and change the `basePath` variable based on your machine/VPS

4. Build the program on your local machine by:
    ```bash
    go mod tidy
    go build
    ```
5. Run the program by running `./arxiv` in the project directory at you terminal

6. You can run the program by running `run.sh` file, simply give the path to the compiled `arxiv` file to it. it will try to run the app for 5 times.

7. The program will generate an RSS feed containing the latest papers on arXiv that match your keywords. You can subscribe to this feed using your preferred RSS reader, such as [Feedreader](https://feedreader.com) or [Akregator](https://apps.kde.org/akregator/).

8. use crontab to run `./arxiv` every 24 hours.
9. Enjoy staying up-to-date on the latest research in your areas of interest!

## Notes:


- arXivSquirrel uses the [RSS news feeds](https://arxiv.org/help/rss) to fetch the papers, please make sure you are following the [arXiv API usage policy](https://arxiv.org/help/api/user-manual)

- Keep in mind that arXivSquirrel is not a substitute for reading the full text of a paper. It's merely a tool to help you identify papers that may be of interest to you.

## Contribution

This is a weekend personal project, It can be way better and optimized, But I;m not going to spend any mode time on it. so sorry it it's not that user friendly!

If you would like to contribute to the development of arXivSquirrel or report any issues, please feel free to open a pull request or issue on this Git

## License

arXivSquirrel is released under the MIT License. See the [LICENSE](https://github.com/mh-salari/arXivSquirrel/blob/master/LICENSE) file for details.

