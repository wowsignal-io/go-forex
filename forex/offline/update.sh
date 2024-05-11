#!/bin/bash

# set -e

ECB_URL='https://www.ecb.europa.eu/stats/eurofxref/eurofxref-hist.zip'
BOC_URL='https://www.bankofcanada.ca/valet/observations/group/FX_RATES_DAILY/csv?start_date=2017-01-03'
RBA_URL='https://www.rba.gov.au/statistics/tables/csv/f11.1-data.csv'

ECB_FILENAME='eurofxref-hist.zip'
BOC_FILENAME='boc_offline_rates.csv'
RBA_FILENAME='rba_offline_rates.csv'

read -r -d '' ECB_LICENSE_HEADER << EOF || true
The file 'eurofxref-hist.zip' was obtained from the website of the European
Central Bank (ECB) on $(date +%Y-%m-%d) using the URL
${ECB_URL}.
EOF

read -r -d '' ECB_LICENSE_TEXT << EOF || true
The rest of this section is an excerpt from ECB's copyright policy retrieved on
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

read -r -d '' BOC_LICENSE_HEADER << EOF || true
The file 'boc_offline_rates.csv' was obtained from the website of the Bank of
Canada (BOC) on $(date +%Y-%m-%d) using the URL
${BOC_URL}.
EOF

read -r -d '' BOC_LICENSE_TEXT << EOF || true
The rest of this section is a copy of BOC's terms of use policy retrieved on
2024-05-11 from https://www.bankofcanada.ca/terms/.

Terms of Use and Disclaimers
This website is provided by the Bank of Canada as a service to its users. Users are defined for the purposes of this statement as “You”. Your access to, and use of, this website constitutes your agreement to accept these Terms of Use and Disclaimers.

If you do not agree with these Terms of Use, or any part of them, you may leave the website.

The Bank of Canada reserves the right to update or modify these Terms of Use and Disclaimers at any time without prior notice. Your continued use of this website following any changes constitutes your agreement to accept such changes.

See also the terms of use for accessing the Bank's social media.
Terms of Use

1. Copyright / Permission to Reproduce

Unless otherwise stated the copyright and any other rights in the contents of the material available through this website, including any images and text, are owned by the Bank of Canada. 

Content on this website is produced and/or compiled by the Bank of Canada for the purpose of providing users with information related to the activities of the Bank. Except as indicated below (section 2, Exceptions to Permission to Reproduce), the Bank permits you to freely use, copy, distribute and transmit its website content under the following terms:

1.1. Attribution

You must attribute the Bank of Canada as the source of the content, and indicate if changes were made.  You may do so in any reasonable manner, but not in any way that suggests that the Bank endorses you or your use of the content.

1.2 Accuracy

You must exercise due diligence in ensuring the accuracy of any content you reproduce.

1.3 Use of Content: Notice Requirement

If You provide content from this website through paid services or incorporate any content in documents for sale (regardless of the medium), You must inform any prospective purchaser, prior to its distribution or sale, that said content was obtained from this website and that such information is available on this website free of charge.

2. Exceptions to Permission to Reproduce

2.1. Bank Note Images

Permission to reproduce bank note images must be obtained in writing from the Bank. See the Bank of Canada Policy on the Reproduction of Bank Note Images for further information.

2.2. Reproduction of the Bank’s Logo and Wordmark

The Bank’s logo and wordmark may not be reproduced, whether for personal, commercial or non-commercial purposes, without written authorization from the Bank. Permission to use the Bank’s Logo or Watermark may be obtained by submitting a request in writing by email.

2.3 Third-party Content and Licensed Content

Some website content may be produced by third parties who retain all rights in such material. Such content is used by the Bank with the permission of the content’s owner and may be subject to specific terms and conditions. The rights of such third parties is protected under the Copyright Act and international agreements.  The Bank does not grant You permission to re-use such third-party content unless it expressly indicates that You may do so.

The Bank will endeavor to acknowledge third party copyright and restrictions on the use of specific content. If third party copyright or restrictions have not been properly acknowledged, please notify us by email so that we can make the appropriate arrangements.

3. No Unlawful or Prohibited Use

As a condition of your use of this website and its contents (the “Site”), you agree to not: 

Use the Site for any purpose that is (i) a violation of any applicable statute, or otherwise unlawful, or (ii) prohibited by these Terms of Use. 
Circumvent such limit or limits imposed by the Bank with respect to the number or frequency of requests made to a Bank Site, including the retrieval of financial data and information using Bank of Canada services (e.g., the Bank of Canada Valet API).
Access or use the Site in a manner that has or could have the effect of disabling, damaging, overburdening, interfering with or impairing the functioning of the Site.
Interfere with, or cause any interference with, any other party's use and enjoyment of the Site.
Engage in any conduct, as determined solely by the Bank, that may harm the Bank or any user of the Site or expose them to liability. 
Knowingly transmit viruses or disabling features which may or will damage, or interfere with, the Site. 
Facilitate or assist any other person in any unlawful or prohibited use.
4. Official Languages

The Bank of Canada observes the Official Languages Act, and is committed to ensuring information and services on this site are available in both English and French. Please note that, with respect to research papers, such papers are generally published in the official language in which they were originally written. In such instances an abstract of the paper is provided in the other official language.

5. Privacy Notice

Personal information collected by the Bank of Canada is protected under the Privacy Act. For further information about the Bank of Canada’s privacy practices, please consult the Bank’s General Privacy Policy.

Please note that listings of job opportunities currently available at the Bank of Canada are contained on a third-party site. Any personal information submitted through an application for employment using that site is protected under the Privacy Act.

Disclaimers

General

The information on this website is provided for general reference purposes only.  While every effort is made to ensure that the site is up to date and accurate, the Bank of Canada accepts no responsibility or liability for the accuracy or completeness of the content or for any loss which may arise from reliance on information contained in this website.

You acknowledge that your use of the information and data, including rates and statistics, provided through this website is at your sole and own risk. Under no circumstances shall the Bank of Canada, its employees, directors, officers, agents, vendors or licensors involved in creating, producing or delivering this site or content found on this site be liable for any loss, injury, claim, liability or damage of any kind resulting in any way from: (a) any content or material downloaded from this website; (b) any links provided on this website; (c) any errors in, or omissions from, the data found on this website; (d) the unavailability or delay of any such data; or (e) your use of the data found on this website or any conclusions you draw from it, regardless of whether you received any assistance from the Bank of Canada or its employees with regard to such data. Under no circumstances is the Bank of Canada liable to you for any amount.

The Bank of Canada provides no warranty, express or implied, as to the accuracy, timeliness, completeness, merchantability, fitness for any particular purpose, title, quality or non-infringement of any service or information contained on the website in any form or manner whatsoever.

Links to other sites are provided for your convenience. The Bank of Canada accepts no responsibility or liability for the content of those sites or of any external site which links to this site.

Rates and Statistics

While statistical data on the Bank of Canada website is derived from sources that the Bank generally considers reliable, the Bank cannot guarantee the completeness or accuracy of such data.

Exchange Rates

All Bank of Canada exchange rates are indicative rates only, derived from averages of transaction prices and price quotes from financial institutions. As such, they are intended to be broadly indicative of market prices at the time of publication but do not necessarily reflect the rates at which actual market transactions have been or could be conducted. They may differ from the rates provided by financial institutions and other market sources. Bank of Canada exchange rates are released for statistical and analytical purposes only and are not intended to be used as the benchmark rate for executing foreign exchange trades. The Bank of Canada does not guarantee the accuracy or completeness of these exchange rates. The underlying data is sourced from Refinitiv (formerly Thomson Reuters).
EOF

read -r -d '' RBA_LICENSE_HEADER << EOF || true
The file 'rba_offline_rates.csv' was obtained from the website of the Reserve Bank
of Australia (RBA) on $(date +%Y-%m-%d) using the URL
${RBA_URL}.
EOF

read -r -d '' RBA_LICENSE_TEXT << EOF || true
The rest of this section is a copy of RBA's terms of use policy retrieved on
2024-05-11 from https://www.rba.gov.au/copyright/.

Copyright and Disclaimer Notice

On This Page
1. Summary
Use of material published by the Reserve Bank of Australia (RBA) on this website, on an RBA app and/or on market data services is subject to the terms and conditions set out in this Notice. Here is a summary of some key elements:
the RBA Logo may be used only in accordance with the Logo Use Guidelines
images of Banknotes may be used only in accordance with the RBA's guidelines for Reproducing Banknotes
content obtained from a third party may be reproduced, published, communicated to the public or adapted only with the permission of the third party. In particular, this website contains data sourced from the Australian Bureau of Statistics (ABS). The conditions for using ABS data are available on the ABS website
the Cash Rate, the Chart Pack, other financial data and econometric models are subject to special conditions as set out in Sections 4, 5 and 6 below
photographs in the Image Library may be reproduced and published without alteration in a context that is not inappropriate
multimedia (not covered by one of the categories above) may be used only for educational purposes, but may not be adapted or published and, in the case of any music, may not be used as a stand-alone file
most other material published on this website is provided under a Creative Commons Attribution 4.0 International License, and can be downloaded, reproduced, published and communicated to the public provided that the RBA is properly attributed.
This is only a summary and it is important that users read this Notice in full to ensure that any proposed use is permitted.
2. Copyright
© Reserve Bank of Australia
Apart from any use as permitted under the Copyright Act 1968, and the permissions explicitly granted below, all other rights are reserved.
The RBA publishes materials (RBA Material) on this website and/or on market data services such as Refinitiv and Bloomberg (Market Data Services) and/or or on any application (App) developed by or on behalf of the RBA available to the users of mobile phones and tablet devices. With the exception of Third Party Material (as defined below), all RBA Material, including (but not limited to) the Excluded Material (as defined below), is the copyright of the RBA.
With exception of the Excluded Material, all RBA Material is provided under a Creative Commons Attribution 4.0 International License (CC BY 4.0 Licence) and may be used in accordance with the terms of that licence. The materials covered by this licence may be reproduced, published, communicated to the public and adapted provided that the RBA is properly attributed in accordance with Section 3 below. Use of these materials is also subject to the disclaimers outlined in Section 7 below.
The terms and conditions of the CC BY 4.0 Licence, as well as further information regarding the licence, can be accessed at <https://creativecommons.org/licenses/by/4.0/legalcode>.
Creative Commons License
This work is licensed under a Creative Commons Attribution 4.0 International License.
The following RBA Material is not provided under the CC BY 4.0 Licence and may be used only in accordance with the following permissions (Excluded Material):
RBA logo: The RBA logo may be used only in accordance with the Logo Use Guidelines.
Banknotes: Images and partial images of past and present Australian banknotes may be used only in accordance with the RBA's guidelines for Reproducing Banknotes.
Third Party Material: Material containing, or derived from or prepared using, content obtained from a third party (Third Party Material), whether it has been:
reproduced in whole or in part in the form supplied by the third party and published as having a third party source; or
used by the RBA to derive or prepare other material published as having both a RBA and a third party source,
may not be reproduced, published, communicated to the public, adapted or otherwise used in whole or part without obtaining the consent of the third party (specifically or, if the material has been published by a third party, under published licence terms or other terms of use). This includes, but is not limited to, data, graphs and tables that show a third party source, alone or with a RBA source, conference and workshop papers authored by third parties, and third party submissions made to the RBA in the course of consultation with the broader community and agencies of government on matters relating to the RBA's responsibilities.
Cash Rate: Use of the Interbank Overnight Cash Rate (the Cash Rate) and materials (including data, graphs, tables, webpages or other publications made available by the RBA) that include the Cash Rate (Cash Rate Materials), is permitted subject to the terms and conditions specified in Section 4 below.
RBA Financial Data: Use of information, rates, facts or knowledge represented in numerical or statistical form (including any annotations) published by the RBA other than the Cash Rate (Financial Data) and the use of materials (including graphs, tables, webpages or other publications made available by the RBA) that include Financial Data (Financial Data Materials) is permitted subject to the terms and conditions specified in Section 5 below.
Econometric models/code: Use of econometric models or code published by the RBA on this website is permitted subject to the terms and conditions specified in Section 6 below.
Linking to the RBA website: Linking to this website is permitted subject to the terms and conditions specified in Section 8 below.
Multimedia: Any video footage, webcast or audio recording, and any photograph, design, icon, logo or other still image or visual representation (that does not fall into any of the other categories of Excluded Material) published on this website or an RBA App or accessed via an RBA social media platform linked to this website may not be reproduced, published, communicated to the public (other than via a link to the relevant page of this website), adapted or otherwise used in whole or part. The only exceptions to this prohibition are:
the photographs in the Image Library on this website, which may be reproduced, published and communicated to the public without alteration or distortion in a context that is not inappropriate, derogatory or offensive; and
any material in this category h. other than the photographs in the Image Library may be used (including by download and reproduction) for educational purposes, but:
may not be published, distributed to the general public or adapted and may not be used for any purpose (commercial or otherwise) that is not an educational one; and
in the case of any music in a video, webcast or audio recording, may not be downloaded, republished, reproduced or otherwise used as a stand-alone file for any purpose including an educational one.
Other material expressly excluded: Any material on which it is indicated, either expressly or implicitly, that the CC BY 4.0 Licence does not apply may not be reproduced, published, communicated to the public, adapted or otherwise used in whole or part, unless the material is stated to be subject to other specified permissions, in which case the use of this material is subject to the terms of those other specified permissions.
Site infrastructure: Any scripts, styles, style sheets or fonts that relate to the structure or format of this website or an RBA App or the visual presentation of text or images on this website or an RBA App (rather than the text or images themselves) may not be reproduced, published, communicated to the public, adapted or otherwise used in whole or part.
Apart from the permissions explicitly granted above and any use as permitted under the Copyright Act 1968, all other rights in respect of the Excluded Material are reserved.
3. Attribution of RBA
Use of RBA Material, whether under the CC BY 4.0 Licence or otherwise, requires you to attribute the work in the manner specified by the RBA. Attribution cannot be done in any way that suggests that the RBA endorses you or your use of the RBA Material.
Unless any RBA Material specifies otherwise, the following form of attribution of RBA Material is required:
Source: Reserve Bank of Australia [year] OR Source: RBA [year]
For RBA Material with identified authors (such as Bulletin articles or speeches) it is acceptable to refer to the authors when referencing, as long as the RBA Material is attributed to the RBA as part of the full reference included in the reference list or bibliography.
4. Cash Rate terms and conditions
The RBA is the administrator of the Cash Rate, which is administered in accordance with the Cash Rate Procedures Manual. The following terms and conditions govern use of the Cash Rate and the Cash Rate Materials.
The Cash Rate and the Cash Rate Materials may be used, reproduced, published, communicated to the public or otherwise referenced for personal or commercial use only if it is not stated, represented or in any way implied (other than in respect of proper attribution as required by Section 3 above) that the RBA endorses any use, reproduction, publication, communication to the public or referencing of the Cash Rate or Cash Rate Materials, or any product, service or financial instrument that is created from, relies on, references or is otherwise derived from the Cash Rate and/or the Cash Rate Materials.
Users of the Cash Rate and/or the Cash Rate Materials are prohibited from using the Cash Rate and/or the Cash Rate Materials for any unlawful purpose. Users of the Cash Rate and/or the Cash Rate Materials agree to refrain from using the Cash Rate and/or the Cash Rate Materials in any way that violates any applicable law or regulation in force in Australia or a foreign country (or any part of Australia or a foreign country).
Users of the Cash Rate and/or the Cash Rate Materials must not engage in improper commercial exploitation of the Cash Rate or the Cash Rate Materials. Improper commercial exploitation includes but is not limited to:
use of the Cash Rate and/or the Cash Rate Materials dishonestly to obtain a benefit, or cause a loss, by deception or other means; and/or
misrepresenting the Cash Rate and/or the Cash Rate Materials as being attributable to, or deriving from some source other than the RBA; and/or
charging a fee to customers for access to the Cash Rate and/or Cash Rate Materials without informing customers that the Cash Rate and/or Cash Rate Materials are published on this website without a fee being charged by the RBA.
Use of the Cash Rate and the Cash Rate Materials is also subject to the disclaimers set out in Section 7 below.
Apart from the specific permissions referred to above, the RBA otherwise reserves and maintains all rights in respect of the Cash Rate and the Cash Rate Materials.
5. Financial Data terms and conditions
From time to time the RBA compiles, prepares or publishes Financial Data, which is not formally administered by the RBA as a financial benchmark. The following terms and conditions govern the use of Financial Data and Financial Data Materials.
Financial Data and Financial Data Materials are made available by the RBA on the understanding that the RBA does not administer the Financial Data or the Financial Data Materials as a benchmark. For this reason, the RBA does not recommend use of any Financial Data or Financial Data Materials as a Financial Benchmark or for any other particular purpose. Users of any Financial Data and/or Financial Data Materials assume the entire risk related to their use of any Financial Data and/or Financial Data Materials, including the use of any materials as the basis for a financial instrument or transaction or any other commercial activity. The RBA does not endorse or promote any financial instrument, transaction or other use (be that commercial or non-commercial) that references or relies on the Financial Data and/or the Financial Data Materials.
Subject to the permissions outlined above and below, the Financial Data and Financial Data Materials may be used, reproduced, published, communicated to the public or otherwise referenced for personal or commercial use only if it is not stated, represented or in any way implied (other than in respect of proper attribution as required by Section 3 above) that the RBA endorses any use, reproduction, publication, communication to the public or referencing of the Financial Data and/or Financial Data Materials or any product, service or financial instrument that is created, relies on, references or is otherwise derived from the Financial Data and/or Financial Data Materials.
Additionally, certain Financial Data and/or Financial Data Materials may contain, derive from or have been prepared using content obtained from a third party. Such material may not be reproduced, published, communicated to the public, adapted, referenced or otherwise used without obtaining the consent of the third party (specifically or, if the material has been published by a third party, under published licence terms or other terms of use). The RBA accepts no responsibility for the unauthorised use of Third Party Material.
The RBA may withdraw, modify or amend any Financial Data and/or Financial Data Materials that appear on this website and/or on the Market Data Services and may alter the methods of calculation, publication schedule, methodology or availability of the Financial Data and/or Financial Data Materials at any time and without notice. The RBA is not, under any circumstances, liable to any user for damages of any kind arising out of or in connection with any such alteration of the Financial Data and/or Financial Data Materials.
Users of the Financial Data and/or Financial Data Materials are prohibited from using the Financial Data and/or Financial Data Materials for any unlawful purpose. Users of the Financial Data and/or Financial Data Materials agree to refrain from using the Financial Data and/or Financial Data Materials in any way that violates any applicable law or regulation in force in Australia or a foreign country (or any part of Australia or a foreign country).
Users of the Financial Data and/or Financial Data Materials must not make improper commercial exploitation of the Financial Data or Financial Data Materials. Improper commercial exploitation includes but is not limited to:
use of Financial Data and/or Financial Data Materials to dishonestly obtain a benefit, or cause a loss, by deception or other means; and/or
misrepresenting the Financial Data and/or Financial Data Materials as being attributable to, or deriving from some other source (except, where the Financial Data and/or Financial Data Materials contain material derived from, or have been prepared using content obtained from, a third party, that third party); and/or
charging a fee to customers for access to the Financial Data and/or Financial Data Materials without informing customers that the Financial Data and/or Financial Data Materials are published on this website without a fee being charged by the RBA.
Use of the Financial Data and Financial Data Materials is also subject to the disclaimers set out in Section 7 below.
Apart from the specific permissions referred to above, the RBA otherwise reserves and maintains all rights in respect of the Financial Data and Financial Data Materials.
6. Econometric models/code terms and conditions
The RBA publishes on this website econometric models or code that can be used to replicate the RBA's research results. Use of such econometric models or code is subject to the following terms and conditions of use.
The use, reproduction, publication, communication to the public, adaptation and referencing of the RBA's econometric models or code is allowed on the condition that users give proper attribution in accordance with Section 3 above. The RBA's econometric models or code may be used, reproduced, published, communicated to the public, adapted, or otherwise referenced only if it is not stated, represented or in any way implied (other than in respect of proper attribution as required by Section 3 above) that the RBA endorses any use, reproduction, publication, communication to the public, adaptation or referencing of the econometric models or code, or any product, service or financial instrument that is created, relies on, references or is otherwise derived from the econometric models or code. Neither the name of the RBA nor the names of any of the authors may be used in any way that suggests that the RBA or the authors endorse or promote works derived from any of the econometric models or code without prior written permission.
Users must not engage in improper commercial exploitation of any econometric models or code. Improper commercial exploitation includes but is not limited to:
use of econometric models or code dishonestly to obtain a benefit, or causing a loss, by deception or other means; and/or
misrepresenting the RBA's econometric models or code as being attributable to, or deriving from some source other than the RBA; and/or
charging a fee to customers for access to the econometric model or code without informing customers that the econometric model or code is published on this website without a fee being charged by the RBA.
Users of the econometric models or code are additionally prohibited from using the econometric models or code for any unlawful purpose. Users of the econometric models or econometric code agree to refrain from using the econometric models or code in any way that violates any applicable law or regulation in force in Australia or a foreign country (or any part of Australia or a foreign country).
Use of econometric models or code is also subject to the disclaimers set out in Section 7 below.
Apart from the specific permissions referred to above, the RBA otherwise reserves and maintains all rights in respect of the econometric models or econometric code.
7. Disclaimers
7.1 By the RBA

RBA Material is intended as a general reference for users. It is made available on the understanding that the RBA, as a result of providing this information, is not engaged in providing professional or financial advice.
While the RBA will make every effort to maintain up-to-date and accurate information on this website, an RBA App and/or on the Market Data Services, users should be aware that the RBA accepts no responsibility for the accuracy or completeness of any RBA Material and recommends that users exercise their own care and judgment with respect to its use.
Users of RBA Material assume the entire risk related to their use of such materials, including the use of any materials as the basis for a financial instrument, transaction or any other commercial activity. The RBA does not accept any liability arising from reliance on or use of any RBA Material.
The RBA does not endorse or promote any financial instrument, transaction or other use (be that commercial or non-commercial) that references or relies on any RBA Material. The RBA expressly disavows any use of RBA Material that in any way violates any applicable law or regulation in force in Australia or a foreign country (or any part of Australia or a foreign country).
The RBA is not, under any circumstances, liable for damages of any kind arising out of or in connection with use of or inability to use any RBA Material, including damages arising from negligence on the part of the RBA, its employees or agents. By using RBA Material, the user agrees to waive all claims against the RBA and its officers, agents, and employees from any and all liability for claims, damages, costs and expenses of any kind arising from or in any way connected to use of any RBA Material, including claims arising from negligence on the part of the RBA, its employees or agents.
Any use of materials provided under the CC BY 4.0 Licence are additionally subject to the disclaimers and warranties as set out in that licence. The terms and conditions can be accessed at: <https://creativecommons.org/licenses/by/4.0/legalcode>.
The RBA has made all reasonable efforts clearly to label material as having a third party source when that material contains, has been derived from or prepared using content obtained from a third party. The RBA uses Third Party Material with permission under licence, and does not guarantee or warrant that the content or information derived from it is accurate, complete or up to date.
This website contains links to the websites of other organisations both within and outside Australia. An RBA App may also contain links to the websites of other organisations. These links to other websites are provided to help meet your needs but the RBA does not provide any warranty in relation to, or take any responsibility for, the content or any other aspect of those websites or of any site or app store through which an RBA App can be downloaded.
The listing of a person or organisation in any part of this website or an RBA App in no way implies any form of endorsement by the RBA of the products or services provided by that person or organisation.
7.2 Third party copyright and disclaimers

Australian Bureau of Statistics (ABS) data are used with permission; copyright for ABS data resides with the Commonwealth of Australia. The conditions for using ABS data are available on the ABS website.
IHS Global Pte Ltd (IHS) disclaims any and all liability associated with any proprietary IHS information used on this site. IHS advises readers to use their own judgment and due diligence to verify any such IHS information.
Separate copyright and disclaimer notices are provided relating to data from APM, ASX, BLADE, HILDA Survey, RP Data Pty Ltd trading as CoreLogic Asia Pacific and ICE Data Indices.
Other third parties have rights not mentioned in this Section.
8. Linking to the RBA website
You may link to this website, provided you do so in a way that is fair and legal and does not damage or take advantage of the RBA's reputation.
You must not:
establish a link in such a way as to suggest any form of association, approval or endorsement with or by the RBA where none exists;
establish a link to this website in any website that is not owned by you; or
frame this website on any other site.
The RBA reserves the right to withdraw linking permission without notice.
Enquiries concerning the use of material on this website should be sent to: rbainfo@rba.gov.au
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
curl -L "${BOC_URL}" > "${BOC_FILENAME}"
curl -L "${RBA_URL}" > "${RBA_FILENAME}"

# Update the license:

{
    echo "============= ECB License ============="
    echo "${ECB_LICENSE_HEADER}"
    echo
    echo "${ECB_LICENSE_TEXT}"
    echo
    echo "============= BOC License ============="
    echo
    echo "${BOC_LICENSE_HEADER}"
    echo
    echo "${BOC_LICENSE_TEXT}"
    echo
    echo "============= RBA License ============="
    echo
    echo "${RBA_LICENSE_HEADER}"
    echo
    echo "${RBA_LICENSE_TEXT}"
} > ./LICENSE || die


popd > /dev/null
