# wishlist
![build status](https://github.com/jckimble/wishlist/actions/workflows/build.yml/badge.svg?branch=master)

wishlist is an cli tool for searching wish.com from the commandline with basic filters

---
- [wishlist](#wishlist)
  - [Install](#install)
  - [Usage](#usage)
  - [License](#license)

---

## Install
```sh
go get -u github.com/jckimble/wishlist
```

## Usage
```sh
$ wishlist
wishlist [options] search query
  -help
		Show this message
  -max int
		Maximum price
  -min int
		Minimum price
  -pages int
		Number of pages (25 per page) (default 4)
  -score float
		Required Keyword Score (default 0.75)
```

## License

Copyright 2018 James Kimble

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
