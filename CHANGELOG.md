v2.0.3:
 - Added a new commit message to be validated for merge from master to the branch (@artus.andermann)

v2.0.3:
 - Fix IsValidMessage method adding a new exception for merge from master to the branch (@esequiel.virtuoso)

v2.0.2:
 - Fix IsValidMessage method adding a new exception for merge from master to the branch (@esequiel.virtuoso)

v2.0.1:
 - Fix IsValidMessage method to skip commit lint when the message is a merge from master to a developer branch. (@esequiel.virtuoso)
 - Add chore option. (@esequiel.virtuoso)

v2.0.0:
 - Change semantic-release message pattern to 'type(scope?): message here'. (@esequiel.virtuoso)

v1.0.5
 - Fix isSetNewVersion function logic (@esequiel.virtuoso)

v1.0.4
 - Fix versioning problem (@lucas.oliveira)
 - Changed `gitlab-backup.sh` to keep all backups (@lucas.oliveira)

v1.0.3
 - Fix only numbers tag bug. (@lucas.oliveira)
 - Add version to log. (@lucas.oliveira)
 - Added more recent version on semantic released tag. (@lucas.olivera)
 - Changed gitlab-backup.sh script to automatic set the last backup version. (@lucas.oliveira)
 - Changed gitlab-backup.sg script to persist the backup in `srv/gitlab/backups` file. (@lucas.oliveira)

v1.0.2
 - Fix most recent commit logic bug. (@esequiel.virtuoso)

v1.0.1
 - Fix tag name conversion bug. (@esequiel.virtuoso)

v1.0.0
 - First version of semantic-release. (@esequiel.virtuoso)
