#!/bin/bash

# set -e

BOC_URL='https://www.bankofcanada.ca/valet/observations/group/FX_RATES_DAILY/csv?start_date=2017-01-03'
ECB_URL='https://www.ecb.europa.eu/stats/eurofxref/eurofxref-hist.zip'

ECB_FILENAME='eurofxref-hist.zip'
BOC_FILENAME='boc_rates.csv'

read -r -d '' ECB_LICENSE_HEADER << EOF || true
The file 'eurofxref-hist.zip' was obtained from the website of the European
Central Bank (ECB) on $(date +%Y-%m-%d) using the URL
${ECB_URL}.
EOF

read -r -d '' ECB_LICENSE_TEXT << EOF || true
The rest of this file is an excerpt from ECB's copyright policy retrieved on
2022-02-10 from https://www.ecb.europa.eu/services/disclaimer/html/index.en.html

Copyright

Copyright © for the entire content of this website: European Central Bank,
Frankfurt am Main, Germany.

Subject to the exception below, users of this website may make free use of the
information obtained directly from it subject to the following conditions:

1. When such information is distributed or reproduced, it must appear accurately
and the ECB must be cited as the source.

2. Where the information is incorporated in documents that are sold (regardless
of the medium), the natural or legal person publishing the information must
inform buyers, both before they pay any subscription or fee and each time they
access the information taken from this website, that the information may be
obtained free of charge through this website.

3. If the information is modified by the user (e.g. by seasonal adjustment of
statistical data or calculation of growth rates) this must be stated explicitly.

4. When linking to this website from business sites or for promotional purposes,
this website must load into the browser's entire window (i.e. it must not appear
within another website's frame).

As an exception to the above, any reproduction, publication or reprint, in whole
or in part, of documents that bear the name of their authors, such as ECB
Working Papers and ECB Occasional Papers, in the form of a different publication
(whether printed or produced electronically) is permitted only with the explicit
prior written authorisation of the ECB or the authors.
EOF

function script_dir() {
    local relative_dir
    relative_dir="$(dirname "${BASH_SOURCE}")"
    pushd "${relative_dir}" > /dev/null 2>&1
    pwd
    popd > /dev/null 2>&1
}

function die() {
    [[ -n "${2}" ]] && echo "${2}" >&2
    [[ -n "${1}" ]] && exit "${1}"
    exit 255
}

pushd "$(script_dir)" > /dev/null

# Download the files:

curl -L "${ECB_URL}" > "${ECB_FILENAME}"

# Update the license:

{
    echo "${ECB_LICENSE_HEADER}"
    echo
    echo "${ECB_LICENSE_TEXT}"
} > ./LICENSE || die


popd > /dev/null
