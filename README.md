# GO-TAMBOON ไปทำบุญ

This is a small challenge project 

### CONTENTS

* `data/fng.csv.rot128` - A ROT-128 encrypted CSV file.
* `cipher/rot128.go` - Sample ROT-128 encrypt-/decrypter.

### EXERCISE

Write a GO command-line module that, when given the CSV list, calls the Charge API to
make donations by creating a charge for each row in the file and produce a summary at the
end.

Example:

performing donations...
done.

        total received: THB  210,000.00
  successfully donated: THB  200,000.00
       faulty donation: THB   10,000.00

    average per person: THB      534.23
            top donors: Obi-wan Kenobi
                        Luke Skywalker
                        Kylo Ren
```

**Requirements:**

* Decrypt the file using a simple [ROT-128][2] algorithm.
* Make donations by creating a Charge via the Charge API for each row in the
  decrypted CSV.
* Produce a brief summary at the end.


