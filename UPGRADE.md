# Upgrading

This document describes future changes that will affect your current resources used
inside of OpenShift. Each change contains description of the change and information
when that change will happen.


## Origin 1.0.x / OSE 3.0.x

* Currently all build pods have a label named `build`. This label is being deprecated
  in favor of `openshift.io/build.name` in Origin 1.0.x / OSE 3.1.x at which point both
  labels will be supported. All the newly created builds will have just the new label.
  In Origin 1.y / OSE 3.y the support for the old label (`build`) will be removed entirely.
  See #3502.
