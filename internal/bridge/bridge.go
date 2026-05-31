package bridge

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"os/exec"
	"time"

	"github.com/gorilla/websocket"
)

// ... rest of the file remains the same as previous update ...