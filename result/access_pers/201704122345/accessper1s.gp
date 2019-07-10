accessper1s_1data = "accessper1s_1.dat"
stats accessper1s_1data using 2 name "accessper1s_1stats"
accessper1s_01data = "accessper1s_01.dat"
stats accessper1s_01data using 2 name "accessper1s_01stats"
accessper1s_001data = "accessper1s_001.dat"
stats accessper1s_001data using 2 name "accessper1s_001stats"
accessper1s_0001data = "accessper1s_0001.dat"
stats accessper1s_0001data using 2 name "accessper1s_0001stats"
accessper1s_00001data = "accessper1s_00001.dat"
stats accessper1s_00001data using 2 name "accessper1s_00001stats"
accessper1s_000001data = "accessper1s_000001.dat"
stats accessper1s_000001data using 2 name "accessper1s_000001stats"
accessper1s_0000001data = "accessper1s_0000001.dat"
stats accessper1s_0000001data using 2 name "accessper1s_0000001stats"

set xlabel "Time /s"
set ylabel "TCAM access count"
set tics font "Arial, 10"

set logscale x
#set yr[0:100]
set xrange[1:]

plot accessper1s_1data using 1:2  w l title "time width = 1.0",\
     accessper1s_01data using 1:2  w l title "time width = 0.1",\
     accessper1s_001data using 1:2  w l title "time width = 0.01",\
     accessper1s_0001data using 1:2  w l title "time width = 0.001",\
     accessper1s_00001data using 1:2  w l title "time width = 0.0001",\
     accessper1s_000001data using 1:2  w l title "time width = 0.00001",\
     accessper1s_0000001data using 1:2  w l title "time width = 0.000001",\
 