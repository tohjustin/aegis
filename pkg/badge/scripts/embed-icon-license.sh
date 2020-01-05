# script for embedding Font Awesome Free license into icons
for i in assets/icons/**/*.svg
do
  echo "" >> $i
  echo "<!--" >> $i
  echo "Font Awesome Free 5.12.0 by @fontawesome - https://fontawesome.com" >> $i
  echo "License - https://fontawesome.com/license/free (Icons: CC BY 4.0, Fonts: SIL OFL 1.1, Code: MIT License)" >> $i
  echo "-->" >> $i
done
