@0xed4ff583aff4a72f;  # Generate with `capnp id`

using Package = import "/sandstorm/package.capnp";

const command :Package.Manifest.Command = (
  executablePath = "/legacy-bridge",
  args = ["127.0.0.1:8080", "/shell"]
);

const manifest :Package.Manifest = (
  # Increment this whenever you distribute a new update to users.
  # Sandstorm will know to replace older versions of the app.
  appVersion = 0,

  actions = [(
    # Defines a "new document" action to add to the main menu.
    input = (none = void),
    title = (defaultText = "New Shell"),

    # This command will run when a new instance of the app
    # is created after the user clicks on "New Widget".
    command = .command
  )],

  # This command runs when the user opens a pre-existing instance
  # of the app.
  continueCommand = .command
);