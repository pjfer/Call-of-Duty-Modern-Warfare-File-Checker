# Call of Duty: Modern Warfare File Checker

A simple program to help you verify which files from your Call of Duty Modern Warface installation may be corrupted and should be downloaded again.

## Getting Started

To smoothly run this program, you just need to execute the MWFileChecker.exe.

### Prerequisites

You only need to have the Call of Duty Modern Warface installed.

## Guide for the program's execution

When you open the program, you will be presented with something like the image below.

![Without game folder selected](images/program.png?raw=true)

As you can see, you have 3 fields, where 2 of them are already fulfilled, meaning you only need to locate the game folder.

Also, the only button available for you to press is the "Compare the hashes" button, because if you already have previous hashes inside a file and you only want to compare them, you can do it without having to hash, again, the files.

After locating the game's folder, your program should be like in the image below.

![With game folder selected](images/program2.png?raw=true)

Now, you can hash the files by simply pressing the "Hash the files" button, which, in the end of hashing, will enable the "Save the hashes", so you can save the filenames and their respective hashes into the "results/myMD5Values.txt" file (if the second field wasn't changed).

There's also the possibility to automate the process. You only need to press the "Full Run" button, which will automatically hash the files, save them, compare them and, if they exist, save the filenames that may be corrupted into the "results/faultyFiles.txt" file (if the third field wasn't changed).

## Author

**Pedro Ferreira**
