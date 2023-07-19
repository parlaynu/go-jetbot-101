# Jetbot Tools in Go

This repository contains some tools for driving a jetbot, written in Go. 

* https://www.waveshare.com/wiki/JetBot_AI_Kit

Why Go instead of the usual python and/or ROS? Mainly as a learning experience and to see how much support 
there is for Go in a robotics environment. Also, the concurrency and parallelism of Go seems like it could be
a good fit.

So far there is only a very basic drive application; I'm planning on building on it
as I have time.

## Prerequisites

### Golang

You need the golang compiler installed on the jetbot. The download page is [here](https://go.dev/dl/).
I downloaded this one:

    wget https://go.dev/dl/go1.20.5.linux-arm64.tar.gz

Expand and install on the jetbot at `/usr/local/go`:

    sudo tar xzf go1.20.5.linux-arm64.tar.gz -C /usr/local

Then add the path to your `~/.profile` file:

    if [ -d "/usr/local/go/bin" ] ; then
        PATH="/usr/local/go/bin:$PATH"
    fi

That's all. The dependencies needed to build should be downloaded automatically.

## Tools

Simply run `make` at the top level and the tools will be built and placed in a top-level
`bin` directory.

### Stop

Like it says, it stops the robot. If one of the other applications is stopped abruptly,
for example with ctrl-c, and the robot is still on the move, you can stop it with this.

    ./bin/stop

### Drive

This is the base level drive application. It uses the gamepad controller to drive the
wheels directly.

There are two controller options:

- '0' - left and right joysticks directly control the left and right motors (lr)
- '1' - right joystick controls longitudinal and angular speeds (lw)

To stop cleanly, press the 'X' button.

Full usage is:

    Usage: drive [-x] [-s top-speed] [-c controller]
      -c int
        	select the controller to use (0 - lr, 1 - lw)
      -s int
        	set the top speed scale (1 - 8) (default 2)
      -x	motors are wired the other way around

If you've wired your motors the same way I have, then you can start with simply:

    ./drive

If the controls are the wrong way for you, you've wired your motors the opposite way to me. To fix 
it, run the command like this:

    ./drive -x

The speed setting controls the duty cycle of the PWM signal driving the motor and 
sets the top speed. This starts with a value of 1 (12.5%) and ends at 8 (100%).
The default is 2 (25%) and is a pretty good place to start.

