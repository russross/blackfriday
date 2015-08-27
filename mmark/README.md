# Mmark

Building on xml2rfc v2 I-D:

    ./mmark/mmark -xml2 -page id.md > x.xml && xml2rfc --text x.xml && \
    rm x.xml && mv x.txt mmark2rfc2.txt
