ZipAlign
========

A Go replacement for Android's SDK ZipAlign tool. See
https://developer.android.com/studio/command-line/zipalign.html for more
information.

When signing an APK with jarsigner, the returned zip is not aligned properly.
Android optimizes APK by aligning each file on 4 bytes to load them with mmap.
Running zipalign on an unaligned APK fixes its padding in the headers of each
entry to align it properly.

Installation: `go get github.com/mozilla-services/zipalign`

Usage: `zipalign -a 4 -i infile.zip -o outfile.zip`

Example
-------

The command below takes a signed unaligned apk and fixes its padding. Verbosity
is reduced to only show files that require padding. The aligned APK is then
verified with Android's original zipalign tool to make sure the alignment is
correct per Android standard.

```
$ zipalign -i ~/app-rocket-webkit-release-1816-signed.apk -o /tmp/rocket-aligned.signed.apk -v 2>&1| grep padding

2018/01/11 08:16:59  --- res/drawable-hdpi-v4/common_google_signin_btn_icon_dark_normal_background.9.png: padding 3 bytes
2018/01/11 08:16:59  --- res/drawable-hdpi-v4/design_ic_visibility.png: padding 1 bytes
2018/01/11 08:16:59  --- res/drawable-hdpi-v4/googleg_disabled_color_18.png: padding 3 bytes
2018/01/11 08:16:59  --- res/drawable-hdpi-v4/home_pattern.png: padding 1 bytes
2018/01/11 08:16:59  --- res/drawable-hdpi-v4/ic_notification.png: padding 3 bytes
2018/01/11 08:16:59  --- res/drawable-hdpi-v4/logotype.png: padding 1 bytes
2018/01/11 08:16:59  --- res/drawable-mdpi-v4/common_google_signin_btn_icon_dark_normal_background.9.png: padding 3 bytes
2018/01/11 08:16:59  --- res/drawable-mdpi-v4/design_ic_visibility.png: padding 1 bytes
2018/01/11 08:16:59  --- res/drawable-mdpi-v4/googleg_disabled_color_18.png: padding 3 bytes
2018/01/11 08:16:59  --- res/drawable-mdpi-v4/notification_bg_low_normal.9.png: padding 1 bytes
2018/01/11 08:16:59  --- res/drawable-xhdpi-v4/common_google_signin_btn_icon_dark_normal_background.9.png: padding 3 bytes
2018/01/11 08:16:59  --- res/drawable-xhdpi-v4/design_ic_visibility.png: padding 1 bytes
2018/01/11 08:16:59  --- res/drawable-xhdpi-v4/googleg_disabled_color_18.png: padding 3 bytes
2018/01/11 08:16:59  --- res/drawable-xhdpi-v4/notification_bg_low_normal.9.png: padding 1 bytes
2018/01/11 08:16:59  --- res/drawable-xxhdpi-v4/common_google_signin_btn_icon_dark_normal_background.9.png: padding 3 bytes
2018/01/11 08:16:59  --- res/drawable-xxhdpi-v4/design_ic_visibility.png: padding 1 bytes
2018/01/11 08:16:59  --- res/drawable-xxhdpi-v4/googleg_disabled_color_18.png: padding 3 bytes
2018/01/11 08:16:59  --- res/drawable-xxhdpi-v4/home_pattern.png: padding 1 bytes
2018/01/11 08:16:59  --- res/drawable-xxhdpi-v4/ic_notification.png: padding 3 bytes
2018/01/11 08:16:59  --- res/drawable-xxhdpi-v4/logotype.png: padding 1 bytes

$ /opt/android-sdk/build-tools/27.0.3/zipalign -c -v 4 /tmp/rocket-aligned.signed.apk
...
Verification succesful
```
