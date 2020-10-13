function Find-Folders {
    [Reflection.Assembly]::LoadWithPartialName("System.Windows.Forms") | Out-Null
    [System.Windows.Forms.Application]::EnableVisualStyles()
    $browse = New-Object System.Windows.Forms.FolderBrowserDialog
    $browse.ShowNewFolderButton = $false
    $browse.Description = "Select the game directory"
    [void]$browse.ShowDialog()
    $browse.SelectedPath
    $browse.Dispose()
} Find-Folders