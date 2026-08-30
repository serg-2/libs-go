# libs-go

For imagelib need:

wget https://www.libraw.org/data/LibRaw-0.22.2.tar.gz
tar -xvf LibRaw-0.22.2.tar.gz
cd LibRaw-0.22.2

# 2. Конфигурируем проект, явно указывая флаги компилятора C++
./configure --enable-static --disable-shared CXXFLAGS="-O2" LDFLAGS="-lstdc++"
# 3. Собираем ШТАТНЫМ Makefile ТОЛЬКО саму библиотеку (это создаст правильный libraw.a)
make lib/libraw.la

sudo make install-libLTLIBRARIES install-dataDATA
=====================================================
wget https://github.com/mm2/Little-CMS/releases/download/lcms2.19.1/lcms2-2.19.1.tar.gz
tar -xvf lcms2-2.19.1.tar.gz && cd Little-CMS-lcms2.19.1
./configure --enable-static --disable-shared
make
sudo make install
==============================================
export CGO_CFLAGS="-I/usr/local/include"
export CGO_LDFLAGS="-L/usr/local/lib"
export CGO_ENABLED=1

go build -ldflags "-linkmode external -extldflags '-static -lstdc++ -lgomp -llcms2 -lz -lm -ldl -lpthread'" -o <EXECUTABLE> *.go 

strip <EXECUTABLE>
