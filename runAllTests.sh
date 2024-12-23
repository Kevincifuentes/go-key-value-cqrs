
mods=$(go list -f '{{.Dir}}' -m | xargs)
pwd=$(echo "$PWD")
assets_path=$pwd/assets
mkdir -p $assets_path
final_coverage_file="$assets_path"/final.out
echo "mode: set" > $final_coverage_file
for mod in $mods; do
            # shellcheck disable=SC2164
            coverprofile=$assets_path/out.txt
            echo "Processing module $mod"
            (cd "$mod"; go test ./... -coverpkg=go-key-value-cqrs/... -coverprofile=$coverprofile)
            (cat $coverprofile | grep -v ".gen.go" | grep -v "mode: set" >> "$assets_path"/final.out)
            rm $coverprofile
done
go tool cover -html=$final_coverage_file -o "$assets_path"/coverage.html

