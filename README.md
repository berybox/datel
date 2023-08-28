# Datel - Database Management Tool
Datel is an experimental web-based tool designed for easy manipulation of generic data. With Datel, you can easily add, delete, and manipulate collections and documents (called items in Datel) within your database using an intuitive web interface. Datel is designed to provide a user-friendly and accessible experience on both desktop and mobile devices. The database only supports MongoDB as its backend.
Please keep in mind that Datel is not intended for production deployment and assumes only a small number of users. Therefore, it does not include any means of authentication.

Possible use cases are:
- Personal database (e.g. contact list)
- Managing other projects that use MongoDB as a database

## Installation
Datel contains everything you need inside its binary and does not require installation or any dependencies. Just download the appropriate binary for your system. It is necessary to use the Mongo database, for which the appropriate URI is required. Datel works well with (free) [MongoDB Atlas](https://www.mongodb.com/atlas).

## Getting Started
Getting started with Datel is quick and easy. You only need to provide URI of the MongoDB database you want to use.

To start the server with web-view, run:
`./datel run -u <MongoDB_URI>`

It is also possible to specify the address/port on which datel will run by putting parameter `-a`, e.g.:
`./datel run -u <MongoDB_URI> -a :1234`
to run Datel on port `1234`

Access to the Datel environment is user based. To add an administrator account, run:
`./datel add-admin -u <MongoDB_URI> -s <User name> -u <User ID>`
this creates an admin account, which can then be used to add other (non-admin) users.

Datel does not have an authentication mechanism implemented at this moment. For local use, it is possible to specify a single user who will be logged in all the time by specifying the `-o` parameter.
`./datel run -u <MongoDB_URI> -o <User ID>`
will run the server with `<User ID>` always logged in.

For use with multiple users, it is possible to use the http header `X-UserID`, which is used by default to get the user ID. Authentication can thus be easily set up for example using a server supporting mTLS (e.g. using [Caddy](https://caddyserver.com/) the setup is very simple). Different users can then access the database using different certificates.

## Basic principles
Similar to MongoDB, Datel uses terms like Database, Collection and Document (called Item in Datel). Each user has the ability to create their own collections to which they then add individual Items. Before adding Items to a collection, the user must set the individual Fields of the Collection. Each Item then contains the values of these Fields. Fields are of only two types at the moment - text and number. Adding more types is planned.

## TODO
- [ ] Add pagination to the item list
- [ ] Add items search
- [ ] Add more data types
- [ ] Make offcanvas menu fixed on desktop

## License

Datel is licensed under the [MIT License](./LICENSE.md). You are free to use, modify, and distribute this tool in accordance with the terms of the license.

## Disclaimer

**Datel** is provided "as is" without any warranties, express or implied. The use of this software is at your own risk. The authors and contributors of this project disclaim any and all warranties, including but not limited to, the implied warranties of merchantability and fitness for a particular purpose. 

In no event shall the authors or contributors be liable for any direct, indirect, incidental, special, exemplary, or consequential damages (including, but not limited to, procurement of substitute goods or services; loss of use, data, or profits; or business interruption) however caused and on any theory of liability, whether in contract, strict liability, or tort (including negligence or otherwise) arising in any way out of the use of this software, even if advised of the possibility of such damage.