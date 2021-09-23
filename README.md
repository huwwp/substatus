# substatus

substatus is a tool for quickly checking status codes of a common URL pattern in subdomains.

## Install

```
go get -u github.com/huwwp/substatus
```

## Usage

httprobe takes a list of webservers on stdin, extracts the last level subdomain and appends it to the URL (with and without a slash), and then quickly checks response codes.

```
# cat webservers.txt
https://www.example.com
https://images.example.com
# cat webservers.txt | substatus
https://www.example.com/www [404]
https://www.example.com/www/ [404]
https://images.example.com/images/ [404]
https://images.example.com/images [302]
https://images.example.com [200]
https://www.example.com [200]
```

## Credits

Inspired by [tomnomnom's](https://github.com/tomnomnom) various tools.
