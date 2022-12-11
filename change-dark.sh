mkdir theme-change
mkdir theme-change/light-theme
mkdir theme-change/dark-theme
cp -R css theme-change/light-theme
cp -R templates theme-change/light-theme
cd theme-change
curl -o dark-theme.tar 'https://mathizen.org/dark-theme.tar'
tar -xvf dark-theme.tar
rm dark-theme.tar
cp -R dark-theme/* ../
cd ..
echo "Theme changed to Dark!"