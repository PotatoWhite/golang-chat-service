// Description: log
// must be install logrus (go get -u github.com/sirupsen/logrus)

package log

import (
	"github.com/sirupsen/logrus"
	"os"
)

var (
	log    *logrus.Logger
	Info   = logrus.Info
	Infoln = logrus.Infoln
	Infof  = logrus.Infof

	Debug   = logrus.Debug
	Debugf  = logrus.Debugf
	Debugln = logrus.Debugln

	Error   = logrus.Error
	Errorf  = logrus.Errorf
	Errorln = logrus.Errorln

	Warn   = logrus.Warn
	Warnf  = logrus.Warnf
	Warnln = logrus.Warnln

	Fatal   = logrus.Fatal
	Fatalf  = logrus.Fatalf
	Fatalln = logrus.Fatalln

	Panic   = logrus.Panic
	Panicf  = logrus.Panicf
	Panicln = logrus.Panicln

	Trace   = logrus.Trace
	Tracef  = logrus.Tracef
	Traceln = logrus.Traceln

	WithField  = logrus.WithField
	WithFields = logrus.WithFields
	WithTime   = logrus.WithTime

	// SetLevel sets the log level
	SetLevel = logrus.SetLevel
	// GetLevel returns the log level
	GetLevel = logrus.GetLevel
	// AddHook adds a hook to the log hooks
	AddHook = logrus.AddHook
)

// Init initializes the log
func init() {
	log = logrus.New()
	log.SetOutput(os.Stdout)
	log.SetLevel(logrus.DebugLevel)
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	Debug = log.Debug
	Debugf = log.Debugf
	Debugln = log.Debugln

	Info = log.Info
	Infoln = log.Infoln
	Infof = log.Infof

	Warn = log.Warn
	Warnf = log.Warnf
	Warnln = log.Warnln

	Error = log.Error
	Errorf = log.Errorf
	Errorln = log.Errorln

	Fatal = log.Fatal
	Fatalf = log.Fatalf
	Fatalln = log.Fatalln

	Panic = log.Panic
	Panicf = log.Panicf
	Panicln = log.Panicln

	Trace = log.Trace
	Tracef = log.Tracef
	Traceln = log.Traceln

	WithField = log.WithField
	WithFields = log.WithFields
	WithTime = log.WithTime

	SetLevel = log.SetLevel
	GetLevel = log.GetLevel
	AddHook = log.AddHook
}
