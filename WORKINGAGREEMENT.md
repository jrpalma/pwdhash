# Working Agreement

## Requirements
It is assumed that some requirements are fixed and cannot be changed.
For example, returning a JSON object can be dimmed as fixed requirement
because another team might depend on the JSON object.
However, there are some areas where the requirements can be enhanced without
changing its final outcome. Such an example can be a versioned URL path.
Using a versioned URL path does not change the shape of the data, but it helps
with future API changes.

## Enhancements
There are times when the requirements need to be enhanced so that future changes can
be addressed without too much refactoring. In such cases, requirements will be
enhanced to address future needs. These enhancements will be documented to provide 
contextual and historical data for future reference.

## Assumptions
There are times when the requirements are not clear. Engineers will try to make the
reasonable assumptions in order to fulfil the requirements. These assumptions should
always try to provide the best product experience while addressing future product needs
such as maintainability, reliability, and ease of use. 
