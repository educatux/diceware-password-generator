# Testing Environment

## Prerequisites
- Vagrant
- VirtualBox
- Git

## Setting Up Test Environment

1. Clone the repository:
```bash
git clone [repository-url]
cd diceware-generator
```

2. Start test environments:
```bash
vagrant up
```

This will create three test environments:
- Ubuntu 20.04
- Fedora 35
- Debian 11

## Running Tests

1. Connect to a test environment:
```bash
vagrant ssh ubuntu
```

2. Run the test suite:
```bash
cd /vagrant
./test.sh
```

## Test Scenarios

1. Build Tests
- Container building
- Go compilation
- Dependency resolution

2. Functional Tests
- Wordlist loading
- Passphrase generation
- QR code generation
- Special character handling

3. Cross-Platform Tests
- Different Linux distributions
- Different terminal emulators

## Continuous Integration

GitHub Actions will automatically:
1. Build the project
2. Run tests
3. Perform security scan
4. Create release (on tags)

## Manual Testing

To test specific features:
```bash
# Build and run
make build
make run

# Clean and rebuild
make clean
make build

# Run in specific environment
vagrant ssh fedora -c "cd /vagrant && make test"
```

## Security Testing

The CI pipeline includes:
- gosec security scanner
- Dependency checking
- Container scanning

## Troubleshooting

1. Vagrant Issues:
```bash
vagrant destroy -f
vagrant up
```

2. Podman Issues:
```bash
podman system reset
make clean
make build
```

3. Test Failures:
- Check logs in /vagrant/logs
- Verify environment variables
- Check file permissions
