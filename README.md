# Digest
Your changes to your google docs notes in your mailbox.

# Installation
```
go get -v github.com/farnasirim/digest/cmd/digest
go install github.com/farnasirim/digest/cmd/digest
```

## Usage
Have a look at
```
digest --help
```
First time usage:
```
mkdir -p ~/.digest/auth/
cp $PATH_TO_GOOGLE_DRIVE_CREDENTIALS ~/.digest/auth/
 digest --folder=$GOOGLE_DRIVE_FOLDER_NAME \
--smtp-user=you@domain.com \
--smtp-pass=yourpass \
--smtp-server-host=smtp.domain.com \
--persist-confs
```
Subsequent usages:
```
digest
```

Sample output:
```
Looking under folder "subjects" with id "folder-id"
Successfully written "Theoretical Computer Science" "doc-id"
Successfully written "Networks" "doc-id"
Successfully written "Scientific Papers" "doc-id"
Email sent successfully
```

Sample email content: <br>

<div style=" border: 1px solid black ; padding: 25px; ">

<h2> Machine Learning, Probability and Statistics.txt </h2>
<div style="word-wrap: break-word; width:700px; font-family: monospace;">
<font size="4px" color="darkgreen">
 <br> With c not too large, it is robust against outliers.
 <br> Geometrical interpretation: Make inner product theta . x to be >= 1
 <br> On the other hand, we’re minimizing size of the theta with the regularization term B.
 <br> That means we’re demanding the normal to the separator (theta) and an x must be in the same direction for true examples and opposite direction for false examples.
</font>
</div>
<br><br>


<h2> Algorithms and problems.txt </h2>
<div style="word-wrap: break-word; width:700px; font-family: monospace;">
<font size="4px" color="darkgreen">

<h4>Segment Tree </h4>
Update the children with “toProp”, not the ”value”. 
Update self from the children, not the intermediate query values in the get function: 
They not necessarily represent your children since they at least one of your children is only partially included in them because of the query (query has intersected you without supersetting you) 
</font>
</div>

</div>


or when you haven't had any updates in your docs:<br>


<div style=" border: 1px solid black ;  padding: 25px;">

<font color="darkred">
		No new notes! *LOUD GASP* ⊙▃⊙ 
</font>

</div>

Unfortunately setting up the google docs access token takes a few clicks and keystrokes. See [this](https://blog.farnasirim.ir/2018/12/changes-to-my-notes-in-google-docs-in.html
) blog post for full instructions.

## License
MIT
