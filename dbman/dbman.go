//   Onix Config Manager - Dbman
//   Copyright (c) 2018-2020 by www.gatblau.org
//   Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
//   Contributors to this project, hereby assign copyright in this code to the project,
//   to be licensed under the same terms as the rest of the code.
package main

type DbMan struct {
	// db     *Database
	// config *Config
	ready bool
}

// func (m *DbMan) start() error {
// 	// load the configuration file
// 	if cfg, err := NewConfig(); err == nil {
// 		m.config = cfg
// 	} else {
// 		return err
// 	}
//
// 	// initialises the logger
// 	m.setLogger(m.config.LogLevel)
//
// 	// the backend is now ready to receive http connections
// 	m.ready = true
//
// 	// start the service listener
// 	m.listen()
//
// 	return nil
// }

// func (m *DbMan) listen() {
// 	mux := mux.NewRouter()
// 	mux.Use(m.loggingMiddleware)
//
// 	// registers web handlers
// 	log.Trace().Msg("registering web root / and liveliness probe /live handlers")
// 	mux.HandleFunc("/", m.liveHandler)
// 	mux.HandleFunc("/live", m.liveHandler)
//
// 	log.Trace().Msg("registering readiness probe handler /ready")
// 	mux.HandleFunc("/ready", m.readyHandler)
//
// 	if m.config.Metrics {
// 		// prometheus metrics
// 		log.Trace().Msg("metrics is enabled, registering handler for endpoint /metrics.")
// 		mux.Handle("/metrics", promhttp.Handler())
// 	}
//
// 	// creates an http server listening on the specified TCP port
// 	server := &http.Server{
// 		Addr:         fmt.Sprintf(":%s", m.config.Port),
// 		WriteTimeout: 15 * time.Second,
// 		ReadTimeout:  15 * time.Second,
// 		IdleTimeout:  time.Second * 60,
// 		Handler:      mux,
// 	}
//
// 	// runs the server asynchronously
// 	go func() {
// 		log.Trace().Msgf("dbman listening on :%s", m.config.Port)
// 		if err := server.ListenAndServe(); err != nil {
// 			log.Err(err).Msg("failed to start server")
// 		}
// 	}()
//
// 	// creates a channel to pass a SIGINT (ctrl+C) kernel signal with buffer capacity 1
// 	stop := make(chan os.Signal, 1)
//
// 	// sends any SIGINT signal to the stop channel
// 	signal.Notify(stop, os.Interrupt)
//
// 	// waits for the SIGINT signal to be raised (pkill -2)
// 	<-stop
//
// 	// gets a context with some delay to shutdown
// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
//
// 	// releases resources if main completes before the delay period elapses
// 	defer cancel()
//
// 	// on error shutdown
// 	if err := server.Shutdown(ctx); err != nil {
// 		log.Info().Msg("shutting down Terra")
// 		log.Fatal().Err(err)
// 	}
// }
//
// func (m *DbMan) setLogger(logLevel string) {
// 	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
// 	switch strings.ToLower(logLevel) {
// 	case "info":
// 		zerolog.SetGlobalLevel(zerolog.InfoLevel)
// 	case "debug":
// 		zerolog.SetGlobalLevel(zerolog.DebugLevel)
// 	case "error":
// 		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
// 	case "fatal":
// 		zerolog.SetGlobalLevel(zerolog.FatalLevel)
// 	case "trace":
// 		zerolog.SetGlobalLevel(zerolog.TraceLevel)
// 	}
// }
//
// func (m *DbMan) loggingMiddleware(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		log.Trace().Msgf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
// 		// Call the next handler, which can be another middleware in the chain, or the final handler.
// 		next.ServeHTTP(w, r)
// 	})
// }
//
// func (m *DbMan) authenticate(w http.ResponseWriter, r *http.Request) bool {
// 	// if there is a username and password
// 	if len(m.config.AuthMode) > 0 && strings.ToLower(m.config.AuthMode) == "basic" {
// 		if r.Header.Get("Authorization") == "" {
// 			// if no authorisation header is passed, then it prompts a client browser to authenticate
// 			w.Header().Set("WWW-Authenticate", `Basic realm="oxterra"`)
// 			w.WriteHeader(http.StatusUnauthorized)
// 			log.Trace().Msg("Unauthorised request.")
// 			return false
// 		} else {
// 			// authenticate the request
// 			requiredToken := m.newBasicToken(m.config.Username, m.config.Password)
// 			providedToken := r.Header.Get("Authorization")
// 			// if the authentication fails
// 			if !strings.Contains(providedToken, requiredToken) {
// 				// returns an unauthorised request
// 				w.WriteHeader(http.StatusForbidden)
// 				return false
// 			}
// 		}
// 	}
// 	return true
// }
//
// // creates a new Basic Authentication Token
// func (m *DbMan) newBasicToken(user string, pwd string) string {
// 	return fmt.Sprintf("Basic %s",
// 		base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", user, pwd))))
// }
//
// func (m *DbMan) liveHandler(w http.ResponseWriter, r *http.Request) {
// 	w.WriteHeader(http.StatusOK)
// 	_, _ = w.Write([]byte("OK"))
// }
//
// func (m *DbMan) readyHandler(w http.ResponseWriter, r *http.Request) {
// 	if !m.ready {
// 		log.Warn().Msg("DbMan is not ready")
// 		w.WriteHeader(http.StatusInternalServerError)
// 		_, err := w.Write([]byte("DbMan is not ready"))
// 		if err != nil {
// 			log.Error().Err(err)
// 		}
// 	} else {
// 		w.WriteHeader(http.StatusOK)
// 		_, _ = w.Write([]byte("OK"))
// 	}
// }
//
// // deploys database schemas
// func (m *DbMan) installHandler(w http.ResponseWriter, r *http.Request) {
// 	if !m.ready {
// 		log.Warn().Msg("DbMan is not ready")
// 		w.WriteHeader(http.StatusInternalServerError)
// 		_, err := w.Write([]byte("DbMan is not ready"))
// 		if err != nil {
// 			log.Error().Err(err)
// 		}
// 	} else {
//
// 		w.WriteHeader(http.StatusOK)
// 		_, _ = w.Write([]byte("OK"))
// 	}
// }
