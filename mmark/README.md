# Mmark

Building an xml2rfc v2 I-D, when `id.md` is your I-D:

    ./mmark/mmark -xml2 -page id.md > x.xml && xml2rfc --text x.xml && \
    rm x.xml && mv x.txt id.txt
