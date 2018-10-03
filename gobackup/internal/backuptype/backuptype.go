package backuptype

//BackupType specifies the type of a backup
type BackupType int

const (
	Full        BackupType = 0 //Full backup
	Incremental BackupType = 1 //Incremental backup
)

func (btype BackupType) String() string {
	names := [...]string{
		"full",
		"incremental"}

	if btype < Full || btype > Incremental {
		return "Unkown"
	}
	return names[btype]
}
