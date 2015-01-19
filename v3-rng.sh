# extract the rnc convert to rng
#
# ./mmark/mmark -xml -page mmark2rfc.md > x.xml
# xmllint --relaxng xml2rfcv3.rng x.xml 

ruby -r rexml/document -r open-uri -e "puts REXML::XPath.first(REXML::Document.new(open(%q{http://tools.ietf.org/id/draft-iab-xml2rfcv2-00.xml})),%{//artwork[@type='application/relax-ng-compact-syntax']}).text" > xml2rfcv3.rnc && \
    trang xml2rfcv3.rnc xml2rfcv3.rng && rm xml2rfcv3.rnc
