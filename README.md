<h2>Good things come to those who wait. - Work in Progress!</h2>

<div align="center"><pre>
.▄▄ · ▄ •▄ ▄• ▄▌ ▐ ▄ ▄ •▄ 
         ▐█ ▀. █▌▄▌▪█▪██▌•█▌▐██▌▄▌▪     ^...^
             ▄▀▀▀█▄▐▀▀▄·█▌▐█▌▐█▐▐▌▐▀▀▄·    <_* *_>   
        ▐█▄▪▐█▐█.█▌▐█▄█▌██▐█▌▐█.█▌      \_/
 ▀▀▀▀ ·▀  ▀ ▀▀▀ ▀▀ █▪·▀  ▀
<h3>Experience True Privacy</h3>
</pre></div>

## Welcome to Skunk, a cutting-edge peer-to-peer communication platform that enables secure and private messaging with powerful encryption.

<p align="center">
  <a href="https://scherzma.github.io/">Homepage</a> |
  <a href="https://github.com/scherzma/Skunk/wiki">Documentation</a>
</p>

<p align="center">
  <img src="https://img.shields.io/github/stars/scherzma/skunk?style=social" alt="GitHub stars">
</p>

## Features

🔒 Secure Communication
- End-to-End Encryption: Secure your messages with powerful encryption protocols.
- Secure communication over TOR: Leverage the anonymity network to enhance security.

🛡️ Advanced Security
- Keys: Usage of public / private keys for secure communication.

💬 Chat Management
- Invite: Easily invite others to your chat groups.
- Join/Leave: Seamless integration for joining or leaving chat groups.
- Sync: Keep your chats updated.

## Quick Start Guide
### Prerequisites
- Go (version 1.20)
- Tor (For contributers / version 0.4.6.10)

_! Build from source description (in progress)_

### Installation

1. Clone the Skunk Repository
   
   ```bash
   git clone https://github.com/your-repo/skunk.git
   cd skunk
   ```
   
2. Install Dependencies

   ```bash
   go mod tidy
   ```
   
3. Run all tests

   ```bash
   go test ./...
   ```
   
_! go build... (in progress)_   

***

## Examples _(in progress)_
### Creating and Managing Chats

- Create a Chat

  ```bash
  skunk chat create <chat-name>
  ```
  
- Invite Peers
  
  ```bash
  skunk chat invite <peer-username>
  ```
  
- Join a Chat
  
  ```bash
  skunk chat join <chat-name>
  ```
  
- Leave a Chat
  
  ```bash
  skunk chat leave <chat-name>
  ```

***

## <a href="https://github.com/scherzma/Skunk/wiki/Coding-Guidelines">Contribution Guidelines</a>
We adhere to the following coding standards for our Go projects:

- _Effective Go_: Guidelines for writing clear, idiomatic, and efficient Go code.
- _Godoc_: Conventions for documenting Go code using comments.
- _Standard Go Project Layout_: Recommended directory structure for Go projects to promote best practices and consistency.

By following these standards, we ensure that our codebase remains consistent, maintainable, and easy to understand for all team members.

***

## License
Skunk is licensed under the **GPL3** License.

**PSA: This is a student project. We do _not_ recommend using this software for sensitive data, as it may not provide the necessary level of security and confidentiality. To handle sensitive data, use alternative software instead.**
