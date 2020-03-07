package app

import (
	"argovue/kube"
	"bufio"
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func (a *App) streamPodLogs(w http.ResponseWriter, r *http.Request, name, namespace, container string) *appError {
	flusher, ok := w.(http.Flusher)
	if !ok {
		return makeError(http.StatusInternalServerError, "Streaming not supported")
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Transfer-Encoding", "identity")
	w.Header().Set("Connection", "keep-alive")
	w.WriteHeader(http.StatusOK)
	flusher.Flush()

	stream, err := kube.GetPodLogs(name, namespace, container)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("data: %s\n\n", err)))
		flusher.Flush()
		log.Debugf("Error getting pod logs %s/%s/%s, error:%s", namespace, name, container, err)
		<-w.(http.CloseNotifier).CloseNotify()
		return nil
	}
	defer stream.Close()

	reader := bufio.NewReader(stream)
	for {
		str, err := reader.ReadString('\n')
		if err != nil {
			log.Errorf("Log stream read error:%s", err)
			break
		}
		_, err = w.Write([]byte(fmt.Sprintf("data: %s\n\n", str)))
		if err != nil {
			log.Errorf("Log stream write error:%s", err)
			break
		}
		flusher.Flush()
	}
	<-w.(http.CloseNotifier).CloseNotify()
	return nil
}
