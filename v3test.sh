wget http://trac.tools.ietf.org/tools/xml2rfc/trac/export/1750/vocabulary/v3/latest/xml2rfcv3.rng
xmllint --noout --relaxng xml2rfcv3.rng $1
