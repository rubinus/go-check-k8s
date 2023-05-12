FROM rubinus/ansible-check-k8s

#RUN mkdir -p /root/.ssh
#COPY id_rsa /root/.ssh/
#COPY id_rsa.pub /root/.ssh/
#COPY known_hosts /root/.ssh/

COPY go-check-k8s /opt/
COPY demo.yml /opt/
RUN chmod +x /opt/go-check-k8s

WORKDIR /opt/

ENTRYPOINT ["/opt/go-check-k8s"]
