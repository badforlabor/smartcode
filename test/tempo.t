
USTRUCT(Blueprintable, BlueprintType)
struct {{.Name}}
{
	GENERATED_BODY()
public:
    {{range .MemberList -}}
	{{- NIL}}UPROPERTY(VisibleAnywhere, BlueprintReadOnly)
		{{.VariableType}} {{.VariableName}}; {{ $length := len .Comment }}{{ if gt $length 0 }}{{range .Comment}}{{.}} {{end}} {{end}}
    {{end}}
};

template<>
inline bool IsRepEqual<{{.Name}}>(const {{.Name}}& A, const {{.Name}}& B)
{
{{ $length := len .MemberList -}}
{{- if lt $length 3 -}}
    {{- Tab}}return false;
{{- else -}}
    {{- Tab}}return {{NIL -}}
    {{- range $index, $element := .MemberList -}}
        {{-  if eq $index 2 }}IsRepEqual(A.{{$element.VariableName}}, B.{{$element.VariableName}}){{end -}}
        {{-  if gt $index 2  -}}
            {{NewLine}}{{Tab}}{{Tab}}&& IsRepEqual(A.{{$element.VariableName}}, B.{{$element.VariableName}}){{NIL -}}
        {{- end -}}
     {{- end -}};
{{- end -}}
{{- NewLine -}}
}

USTRUCT(Blueprintable, BlueprintType)
struct {{.Name}}List
{
	GENERATED_BODY()
public:
	UPROPERTY(VisibleAnywhere, BlueprintReadOnly)
		TArray<{{.Name}}> DataList;
};

// .h
{
	UPROPERTY(VisibleAnywhere, BlueprintReadOnly)
		bool bValid_{{.Name}} = false;
	UPROPERTY(VisibleAnywhere, BlueprintReadOnly)
		{{.Name}} Local_{{.Name}};
	{{.Name}} OldLocal_{{.Name}};
	FAutonomousBatchRep<{{.Name}}, {{.Name}}List, UPlayerSmoothMoveComponent> BatRep{{.Name}};
	void OnBatchRepDirty(const {{.Name}}& Data);
	void OnBatchRep(const {{.Name}}List& Data)
	{
		ServerRepCached_{{.Name}}(Data);
	}
	UFUNCTION(Reliable, NetMulticast)
		void MultiRepCached_{{.Name}}(const {{.Name}}List& ListData);
	UFUNCTION(Reliable, Server, WithValidation)
		void ServerRepCached_{{.Name}}(const {{.Name}}List& ListData);
	UPROPERTY(BlueprintAssignable)
		FRepAction OnRepAction_{{.Name}};
	UFUNCTION(BlueprintCallable)
		void ReqRep_{{.Name}}(const {{.Name}}& Data);
}

// .cpp
bool UPlayerSmoothMoveComponent::ServerRepCached_{{.Name}}_Validate(const {{.Name}}List& ListData)
{
	return true;
}
void UPlayerSmoothMoveComponent::ServerRepCached_{{.Name}}_Implementation(const {{.Name}}List& ListData)
{
	MultiRepCached_{{.Name}}(ListData);
}
void UPlayerSmoothMoveComponent::MultiRepCached_{{.Name}}_Implementation(const {{.Name}}List& ListData)
{
	BatRep{{.Name}}.OnRecv(ListData);
}
void UPlayerSmoothMoveComponent::ReqRep_{{.Name}}(const {{.Name}}& Data)
{
	auto One = Data;
	BatRep{{.Name}}.AddOne(One);
}
void UPlayerSmoothMoveComponent::OnBatchRepDirty(const {{.Name}}& Data)
{
	bValid_{{.Name}} = true;
	Local_{{.Name}} = Data;
	OnRepAction_{{.Name}}.Broadcast();
}
void TestFunction(const {{.Name}}& Data)
{
	Local_{{ToUpper .Name}} = Data;
	OnRep_{{StripWord .Name 1}}.Broadcast();
	bValid_{{ToSnake .Name}} = true;
	bValid_{{ToCamel .Name}} = true;
	bValid_{{ToLowerCamel .Name}} = true;
}


