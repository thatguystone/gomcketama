SPY_JAR = spymemcached-2.10.3.jar

mcketama_java_test.go: SpyKetama.class | $(SPY_JAR)
	java -cp '.:$(SPY_JAR)' SpyKetama > $@
	go fmt

SpyKetama.class: SpyKetama.java | $(SPY_JAR)
	javac -cp $(SPY_JAR) SpyKetama.java

$(SPY_JAR):
	wget https://spymemcached.googlecode.com/files/$(SPY_JAR)

clean:
	rm -f SpyKetama.class
	rm -f $(SPY_JAR)
