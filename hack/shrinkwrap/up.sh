#!/bin/sh

#
# up: bring up the services
#

SCRIPTDIR=$(cd $(dirname "$0") && pwd)

echo "$(tput setaf 2)Booting Lunchpail for arch=$ARCH$(tput sgr0)"

for f in "$SCRIPTDIR"/01-jaas-prereqs1.yml "$SCRIPTDIR"/02-jaas.yml "$SCRIPTDIR"/04-jaas-defaults.yml "$SCRIPTDIR"/05-jaas-default-user.yml
do
    if [ -f "${f%%.yml}.namespace" ]; then ns="-n $(cat "${f%%.yml}.namespace")"; else ns=""; fi
    kubectl apply --server-side -f $f $ns

    if [ "$(basename $f)" = "02-jaas.yml" ]
    then
        if which gum > /dev/null 2>&1
        then
            gum spin --title "$(tput setaf 2)Waiting for controllers to be ready$(tput sgr0)" -- \
              kubectl wait pod -l app.kubernetes.io/name=dlf -n jaas-system --for=condition=ready --timeout=-1s && \
                kubectl wait pod -l app.kubernetes.io/part-of=codeflare.dev -n jaas-system --for=condition=ready --timeout=-1s
        else
            echo "$(tput setaf 2)Waiting for controllers to be ready$(tput sgr0)"
            kubectl wait pod -l app.kubernetes.io/name=dlf -n jaas-system --for=condition=ready --timeout=-1s
            kubectl wait pod -l app.kubernetes.io/part-of=codeflare.dev -n jaas-system --for=condition=ready --timeout=-1s
        fi
    fi
done

# Future: wait for nvidia operators, too
#if [[ "$HAS_NVIDIA" = true ]]; then
#    echo "$(tput setaf 2)Waiting for gpu operator to be ready$(tput sgr0)"
#    $KUBECTL wait pod -l app.kubernetes.io/managed-by=gpu-operator -n $NAMESPACE_SYSTEM --for=condition=ready --timeout=-1s
#fi